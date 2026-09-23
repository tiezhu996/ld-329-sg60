<template>
  <el-drawer
    :model-value="modelValue"
    :title="store.detail?.title ?? '求助详情'"
    size="520px"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <div v-if="store.detail" v-loading="store.detailLoading" class="need-detail">
      <div class="need-detail__head">
        <div class="tag-row">
          <el-tag :type="NEED_STATUS_TAG_TYPES[store.detail.status]" effect="plain">
            {{ NEED_STATUS_LABELS[store.detail.status] }}
          </el-tag>
          <el-tag effect="plain">{{ store.detail.category }}</el-tag>
          <el-tag effect="plain">{{ store.detail.campus }}</el-tag>
          <el-tag type="success" effect="plain">{{ store.detail.budgetType }}</el-tag>
        </div>
        <p>{{ store.detail.description }}</p>
        <small>发布人 {{ store.detail.requester }} · 期望时间 {{ store.detail.expectTime }}</small>
      </div>

      <el-alert :title="queueNotice" :type="queueNoticeType" :closable="false" show-icon />

      <section>
        <h3>候选名单（按提交先后排队）</h3>
        <NeedQueueList
          :entries="store.detail.entries"
          :is-publisher="store.detail.isPublisher"
          :need-status="store.detail.status"
          :has-interviewee="Boolean(store.detail.currentInterviewee)"
          :acting="store.acting"
          @pick="handlePick"
          @withdraw="handleWithdraw"
        />
      </section>

      <section v-if="canRespond">
        <h3>我要响应</h3>
        <NeedResponseForm :key="formKey" :acting="store.acting" @submit="handleSubmit" />
      </section>

      <section v-if="store.detail.isPublisher && store.detail.status === 'open'" class="need-detail__footer">
        <el-popconfirm title="关闭后候选名单停止，确定关闭这条求助吗？" @confirm="handleClose">
          <template #reference>
            <el-button type="danger" plain :loading="store.acting">关闭求助</el-button>
          </template>
        </el-popconfirm>
      </section>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { ElMessage } from 'element-plus';
import {
  NEED_QUEUE_MESSAGES,
  NEED_STATUS_LABELS,
  NEED_STATUS_TAG_TYPES,
} from '../../constants/need.constants';
import { useNeedsStore } from '../../stores/needs.store';
import type { SubmitResponsePayload } from '../../types/domain';
import NeedQueueList from './NeedQueueList.vue';
import NeedResponseForm from './NeedResponseForm.vue';

defineProps<{ modelValue: boolean }>();
const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>();

const store = useNeedsStore();
const formKey = ref(0);

const canRespond = computed(() => {
  const detail = store.detail;
  if (!detail || detail.status !== 'open' || detail.isPublisher) return false;
  return detail.myStatus === 'none' || detail.myStatus === 'withdrawn';
});

const queueNotice = computed(() => {
  const detail = store.detail;
  if (!detail) return '';
  if (detail.status === 'closed') return '求助已关闭，候选名单已停止';
  if (detail.isPublisher) {
    return detail.currentInterviewee
      ? `当前约谈人：${detail.currentInterviewee}，还有 ${detail.waitingCount} 人排队等待`
      : `暂无约谈，可从候选名单中挑一位约谈（${detail.waitingCount} 人排队）`;
  }
  switch (detail.myStatus) {
    case 'waiting':
      return detail.currentInterviewee
        ? `你在候选名单第 ${detail.myPosition} 位，当前约谈人：${detail.currentInterviewee}`
        : `你在候选名单第 ${detail.myPosition} 位，等待发布人安排约谈`;
    case 'interviewing':
      return '你正在约谈中，约谈开始前仍可退出轮候';
    case 'withdrawn':
      return '你已退出轮候，可重新提交响应';
    default:
      return '你还未响应，提交能帮的方式和可上门时段即可排队';
  }
});

const queueNoticeType = computed(() => {
  const detail = store.detail;
  if (!detail || detail.status === 'closed') return 'info';
  if (detail.myStatus === 'interviewing') return 'success';
  if (detail.myStatus === 'waiting' || detail.isPublisher) return 'warning';
  return 'info';
});

async function runAction(action: () => Promise<void>, successMessage: string) {
  try {
    await action();
    ElMessage.success(successMessage);
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '操作失败，请稍后重试');
  }
}

async function handleSubmit(payload: SubmitResponsePayload) {
  await runAction(async () => {
    await store.submit(payload);
    formKey.value += 1;
  }, NEED_QUEUE_MESSAGES.submitSuccess);
}

async function handlePick(responseId: number) {
  await runAction(() => store.pick(responseId), NEED_QUEUE_MESSAGES.pickSuccess);
}

async function handleWithdraw(responseId: number) {
  await runAction(() => store.withdraw(responseId), NEED_QUEUE_MESSAGES.withdrawSuccess);
}

async function handleClose() {
  await runAction(() => store.close(), NEED_QUEUE_MESSAGES.closeSuccess);
}
</script>
