<template>
  <div class="panel">
    <h2>技能求助 · 候选名单</h2>
    <el-table v-loading="store.listLoading" :data="store.list" size="small">
      <el-table-column label="需求" min-width="170">
        <template #default="{ row }">
          <div>{{ row.title }}</div>
          <small>{{ row.requester }} · {{ row.expectTime }} · {{ row.budgetType }}</small>
        </template>
      </el-table-column>
      <el-table-column prop="category" label="类别" width="82" />
      <el-table-column prop="campus" label="校区" width="96" />
      <el-table-column prop="responses" label="响应" width="72" sortable />
      <el-table-column label="当前约谈人" width="104">
        <template #default="{ row }">
          <span v-if="row.currentInterviewee">{{ row.currentInterviewee }}</span>
          <span v-else class="muted">—</span>
        </template>
      </el-table-column>
      <el-table-column label="我的轮候" width="104">
        <template #default="{ row }">
          <el-tag v-if="row.myStatus === 'interviewing'" type="success" effect="plain">约谈中</el-tag>
          <span v-else>{{ myQueueText(row) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="88">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row)" effect="plain">
            {{ statusLabel(row) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column width="100">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row.id)">候选名单</el-button>
        </template>
      </el-table-column>
    </el-table>
    <NeedDetailDrawer v-model="drawerVisible" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { NEED_QUEUE_MESSAGES, NEED_STATUS_LABELS, NEED_STATUS_TAG_TYPES } from '../../constants/need.constants';
import { useNeedsStore } from '../../stores/needs.store';
import type { NeedSummary } from '../../types/domain';
import NeedDetailDrawer from './NeedDetailDrawer.vue';

const store = useNeedsStore();
const drawerVisible = ref(false);

function statusLabel(need: NeedSummary): string {
  return NEED_STATUS_LABELS[need.status];
}

function statusTagType(need: NeedSummary): 'success' | 'info' {
  return NEED_STATUS_TAG_TYPES[need.status];
}

function myQueueText(need: NeedSummary): string {
  if (need.isPublisher) return '我发布的';
  if (need.myStatus === 'waiting') return `第 ${need.myPosition} 位`;
  if (need.myStatus === 'withdrawn') return '已退出';
  return '未响应';
}

async function openDetail(needId: number) {
  try {
    await store.loadDetail(needId);
    drawerVisible.value = true;
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : NEED_QUEUE_MESSAGES.detailFailed);
  }
}

onMounted(async () => {
  try {
    await store.loadList();
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : NEED_QUEUE_MESSAGES.loadFailed);
  }
});
</script>
