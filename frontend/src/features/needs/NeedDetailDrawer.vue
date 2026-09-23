<template>
  <el-drawer
    :model-value="visible"
    size="720px"
    :title="detail ? detail.title : '求助详情'"
    @update:model-value="emit('update:visible', $event)"
  >
    <div v-if="detail" class="need-detail">
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="发布人">{{ detail.requester }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="NEED_STATUS_TAG[detail.status]" size="small">{{ detail.statusLabel }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="类别">{{ detail.category }}</el-descriptions-item>
        <el-descriptions-item label="校区">{{ detail.campus }}</el-descriptions-item>
        <el-descriptions-item label="期望时间">{{ detail.expectTime }}</el-descriptions-item>
        <el-descriptions-item label="回报类型">{{ detail.budgetType }}</el-descriptions-item>
      </el-descriptions>
      <p class="need-detail__desc">{{ detail.description }}</p>

      <el-alert
        v-if="detail.status === 'closed'"
        title="求助已关闭，候选名单停止变动"
        type="info"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="detail.currentInterview"
        :title="`当前约谈：${detail.currentInterview} · 同一时间只安排一位约谈，其余候选人继续等待`"
        type="warning"
        :closable="false"
        show-icon
      />
      <el-alert v-else title="暂无约谈，等待发布人从名单中挑选" type="info" :closable="false" show-icon />

      <el-alert
        v-if="detail.myQueue"
        class="my-queue-alert"
        :title="`我的位置：第 ${detail.myQueue.position} 位 · ${detail.myQueue.statusLabel} — ${detail.myQueue.result}`"
        :type="detail.myQueue.status === 'interviewing' ? 'success' : 'warning'"
        :closable="false"
        show-icon
      />

      <h3>候选名单（按提交先后排队，共 {{ detail.responseCount }} 人）</h3>
      <el-empty v-if="detail.queue.length === 0" description="还没有人响应，来第一个排队吧" :image-size="80" />
      <el-table v-else :data="detail.queue" size="small" :row-class-name="rowClass">
        <el-table-column label="位次" width="64">
          <template #default="{ row }">#{{ row.position }}</template>
        </el-table-column>
        <el-table-column label="响应者" width="120">
          <template #default="{ row }">
            {{ row.responder }}
            <el-tag v-if="row.mine" size="small" type="success" effect="plain">我</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="helpMethod" label="能帮的方式" min-width="180" show-overflow-tooltip />
        <el-table-column label="可上门时段" width="170">
          <template #default="{ row }">
            <el-tag v-for="slot in row.timeSlots" :key="slot" size="small" effect="plain" class="slot-tag">{{ slot }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submittedAt" label="提交时间" width="130" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="RESPONSE_STATUS_TAG[row.status as ResponseStatus]" size="small">{{ row.statusLabel }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="110" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="detail.canPick && row.status === 'waiting'"
              type="primary"
              size="small"
              :loading="pending"
              @click="emit('pick', row.id)"
            >约谈 TA</el-button>
            <el-button
              v-else-if="row.mine && row.status === 'waiting'"
              type="danger"
              size="small"
              plain
              :loading="pending"
              @click="emit('exit')"
            >退出排队</el-button>
            <el-button
              v-else-if="row.mine && row.status === 'interviewing'"
              type="danger"
              size="small"
              plain
              :loading="pending"
              @click="emit('exit')"
            >撤回约谈</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="need-detail__actions">
        <el-button v-if="detail.canRespond" type="primary" :loading="pending" @click="emit('respond')">我要响应</el-button>
        <el-popconfirm
          v-if="detail.canClose"
          title="关闭后候选名单停止变动，确定关闭该求助吗？"
          confirm-button-text="关闭求助"
          cancel-button-text="再想想"
          @confirm="emit('close')"
        >
          <template #reference>
            <el-button type="danger" plain :loading="pending">关闭求助</el-button>
          </template>
        </el-popconfirm>
      </div>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { NEED_STATUS_TAG, RESPONSE_STATUS_TAG } from '../../constants/need.constants';
import type { NeedDetail, NeedResponseEntry, ResponseStatus } from '../../types/domain';

defineProps<{ visible: boolean; detail: NeedDetail | null; pending: boolean }>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  pick: [responseId: number];
  exit: [];
  close: [];
  respond: [];
}>();

function rowClass({ row }: { row: NeedResponseEntry }): string {
  if (row.status === 'interviewing') {
    return 'queue-row--interviewing';
  }
  return row.mine ? 'queue-row--mine' : '';
}
</script>
