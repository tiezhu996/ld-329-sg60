<template>
  <div class="queue-list">
    <el-empty v-if="entries.length === 0" description="还没有人响应，来排第一个" :image-size="72" />
    <article
      v-for="entry in entries"
      :key="entry.responseId"
      class="queue-entry"
      :class="{ 'queue-entry--withdrawn': entry.status === 'withdrawn' }"
    >
      <div class="queue-entry__head">
        <el-tag v-if="entry.status === 'interviewing'" type="success">约谈中</el-tag>
        <el-tag v-else-if="entry.status === 'waiting'" type="warning" effect="plain">
          第 {{ entry.position }} 位
        </el-tag>
        <el-tag v-else type="info" effect="plain">已退出</el-tag>
        <strong>{{ entry.responder }}</strong>
        <el-tag v-if="entry.isViewer" size="small" effect="dark">我</el-tag>
        <small class="queue-entry__time">{{ entry.submittedAt }}</small>
      </div>
      <p>{{ entry.helpOffer }}</p>
      <div class="tag-row">
        <el-tag v-for="slot in entry.visitSlots" :key="slot" size="small" effect="plain">
          {{ slot }}
        </el-tag>
      </div>
      <div v-if="showActions(entry)" class="queue-entry__actions">
        <el-button
          v-if="canPick(entry)"
          size="small"
          type="primary"
          :disabled="acting"
          @click="emit('pick', entry.responseId)"
        >
          约谈 TA
        </el-button>
        <el-popconfirm
          v-if="canWithdraw(entry)"
          title="确定退出轮候吗？"
          @confirm="emit('withdraw', entry.responseId)"
        >
          <template #reference>
            <el-button size="small" type="danger" plain :disabled="acting">退出轮候</el-button>
          </template>
        </el-popconfirm>
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import type { NeedQueueEntry, NeedStatus } from '../../types/domain';

const props = defineProps<{
  entries: NeedQueueEntry[];
  isPublisher: boolean;
  needStatus: NeedStatus;
  hasInterviewee: boolean;
  acting: boolean;
}>();

const emit = defineEmits<{
  pick: [responseId: number];
  withdraw: [responseId: number];
}>();

// 发布人挑约谈对象：仅在没有约谈进行时，对排队中的响应者可用
function canPick(entry: NeedQueueEntry): boolean {
  return props.isPublisher && props.needStatus === 'open' && !props.hasInterviewee && entry.status === 'waiting';
}

// 响应者退出轮候：约谈开始前（排队中）和约谈中都可以退出，退出后由下一位接手
function canWithdraw(entry: NeedQueueEntry): boolean {
  return entry.isViewer && props.needStatus === 'open' && entry.status !== 'withdrawn';
}

function showActions(entry: NeedQueueEntry): boolean {
  return canPick(entry) || canWithdraw(entry);
}
</script>
