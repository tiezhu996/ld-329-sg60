<template>
  <div class="panel">
    <div class="panel-header">
      <h2>{{ NEED_BOARD_TITLE }}</h2>
      <el-button size="small" :loading="loading" @click="loadNeeds">刷新</el-button>
    </div>
    <el-alert v-if="error" :title="error" type="error" show-icon class="board-alert" />
    <el-table v-loading="loading" :data="needs" size="small" @row-click="openDetail">
      <el-table-column label="求助" min-width="170">
        <template #default="{ row }">
          <div>{{ row.title }}</div>
          <small>{{ row.requester }} · {{ row.campus }}</small>
        </template>
      </el-table-column>
      <el-table-column prop="category" label="类别" width="76" />
      <el-table-column label="状态" width="86">
        <template #default="{ row }">
          <el-tag :type="NEED_STATUS_TAG[row.status as NeedStatus]" size="small">{{ row.statusLabel }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="响应 / 排队" width="96">
        <template #default="{ row }">{{ row.responseCount }} / {{ row.waitingCount }}</template>
      </el-table-column>
      <el-table-column label="当前约谈人" width="104">
        <template #default="{ row }">
          <el-tag v-if="row.currentInterview" type="warning" size="small" effect="plain">{{ row.currentInterview }}</el-tag>
          <span v-else class="muted">{{ NO_INTERVIEW_TEXT }}</span>
        </template>
      </el-table-column>
      <el-table-column label="我的轮候" min-width="150">
        <template #default="{ row }">
          <template v-if="row.myQueue">
            <div>第 {{ row.myQueue.position }} 位 · {{ row.myQueue.statusLabel }}</div>
            <small>{{ row.myQueue.result }}</small>
          </template>
          <span v-else class="muted">{{ NOT_JOINED_TEXT }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="96" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" plain @click.stop="openDetail(row)">名单</el-button>
        </template>
      </el-table-column>
    </el-table>

    <NeedDetailDrawer
      v-model:visible="drawerVisible"
      :detail="detail"
      :pending="pending"
      @pick="handlePick"
      @exit="handleExit"
      @close="handleClose"
      @respond="formVisible = true"
    />
    <ResponseFormDialog v-model:visible="formVisible" :pending="pending" @submit="handleSubmitResponse" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import NeedDetailDrawer from './NeedDetailDrawer.vue';
import ResponseFormDialog from './ResponseFormDialog.vue';
import { NEED_BOARD_TITLE, NEED_STATUS_TAG, NOT_JOINED_TEXT, NO_INTERVIEW_TEXT } from '../../constants/need.constants';
import {
  closeNeed,
  exitNeedResponse,
  fetchNeedDetail,
  fetchNeeds,
  pickNeedCandidate,
  submitNeedResponse,
} from '../../services/need.service';
import type { NeedDetail, NeedStatus, NeedSummary, SubmitResponsePayload } from '../../types/domain';

const props = defineProps<{ currentUser: string }>();

const needs = ref<NeedSummary[]>([]);
const detail = ref<NeedDetail | null>(null);
const loading = ref(false);
const pending = ref(false);
const error = ref('');
const drawerVisible = ref(false);
const formVisible = ref(false);

onMounted(loadNeeds);

async function loadNeeds() {
  loading.value = true;
  error.value = '';
  try {
    needs.value = await fetchNeeds(props.currentUser);
  } catch (err) {
    error.value = err instanceof Error ? err.message : '加载求助列表失败';
  } finally {
    loading.value = false;
  }
}

async function openDetail(need: NeedSummary) {
  drawerVisible.value = true;
  try {
    detail.value = await fetchNeedDetail(need.id, props.currentUser);
  } catch (err) {
    notifyError(err);
  }
}

async function runAction(action: () => Promise<NeedDetail>, successText: string) {
  pending.value = true;
  try {
    detail.value = await action();
    ElMessage.success(successText);
    await loadNeeds();
  } catch (err) {
    notifyError(err);
  } finally {
    pending.value = false;
  }
}

function handleSubmitResponse(payload: SubmitResponsePayload) {
  if (!detail.value) {
    return;
  }
  const needId = detail.value.id;
  formVisible.value = false;
  void runAction(
    () => submitNeedResponse(needId, props.currentUser, payload),
    '已加入候选名单，按提交先后排队',
  );
}

function handlePick(responseId: number) {
  if (!detail.value) {
    return;
  }
  const needId = detail.value.id;
  void runAction(() => pickNeedCandidate(needId, props.currentUser, responseId), '已安排约谈，其余候选人继续等待');
}

async function handleExit() {
  if (!detail.value?.myQueue) {
    return;
  }
  const interviewing = detail.value.myQueue.status === 'interviewing';
  const confirmText = interviewing
    ? '撤回后将由排队的下一位自动接手约谈，确定撤回吗？'
    : '确定退出候选名单吗？';
  try {
    await ElMessageBox.confirm(confirmText, '确认操作', { type: 'warning' });
  } catch {
    return;
  }
  const needId = detail.value.id;
  void runAction(
    () => exitNeedResponse(needId, props.currentUser),
    interviewing ? '已撤回约谈，下一位候选人已接手' : '已退出排队',
  );
}

function handleClose() {
  if (!detail.value) {
    return;
  }
  const needId = detail.value.id;
  void runAction(() => closeNeed(needId, props.currentUser), '求助已关闭，候选名单停止变动');
}

function notifyError(err: unknown) {
  ElMessage.error(err instanceof Error ? err.message : '操作失败，请稍后重试');
}
</script>
