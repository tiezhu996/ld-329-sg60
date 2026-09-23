<template>
  <div class="response-form">
    <el-input
      v-model="helpOffer"
      type="textarea"
      :rows="3"
      maxlength="200"
      show-word-limit
      placeholder="描述你能帮的方式，例如：可带设备上门、包修图、可提供课程笔记……"
    />
    <div class="response-form__slots">
      <span class="muted">可上门时段：</span>
      <el-checkbox-group v-model="visitSlots">
        <el-checkbox v-for="slot in VISIT_SLOT_OPTIONS" :key="slot" :value="slot">{{ slot }}</el-checkbox>
      </el-checkbox-group>
    </div>
    <el-button type="primary" :loading="acting" @click="handleSubmit">提交响应并排队</el-button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage } from 'element-plus';
import { NEED_QUEUE_MESSAGES, VISIT_SLOT_OPTIONS } from '../../constants/need.constants';
import { useNeedsStore } from '../../stores/needs.store';
import type { SubmitResponsePayload } from '../../types/domain';

defineProps<{ acting: boolean }>();
const emit = defineEmits<{ submit: [payload: SubmitResponsePayload] }>();

const store = useNeedsStore();
const helpOffer = ref('');
const visitSlots = ref<string[]>([]);

function handleSubmit() {
  if (!helpOffer.value.trim()) {
    ElMessage.warning(NEED_QUEUE_MESSAGES.needHelpOffer);
    return;
  }
  if (visitSlots.value.length === 0) {
    ElMessage.warning(NEED_QUEUE_MESSAGES.needVisitSlots);
    return;
  }
  emit('submit', {
    responder: store.viewer,
    helpOffer: helpOffer.value.trim(),
    visitSlots: [...visitSlots.value],
  });
}
</script>
