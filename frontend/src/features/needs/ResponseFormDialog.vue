<template>
  <el-dialog
    :model-value="visible"
    title="响应求助 · 加入候选名单"
    width="480px"
    @update:model-value="emit('update:visible', $event)"
    @closed="reset"
  >
    <el-form label-position="top">
      <el-form-item label="能帮的方式" required>
        <el-input
          v-model="helpMethod"
          type="textarea"
          :rows="3"
          maxlength="120"
          show-word-limit
          placeholder="说说你打算怎么帮，例如：可带设备上门拍摄并精修 9 张"
        />
      </el-form-item>
      <el-form-item label="可上门时段（至少选一项）" required>
        <el-checkbox-group v-model="timeSlots">
          <el-checkbox v-for="slot in VISIT_TIME_SLOTS" :key="slot" :value="slot">{{ slot }}</el-checkbox>
        </el-checkbox-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:visible', false)">取消</el-button>
      <el-button type="primary" :loading="pending" :disabled="!canSubmit" @click="submit">提交并排队</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { VISIT_TIME_SLOTS } from '../../constants/need.constants';
import type { SubmitResponsePayload } from '../../types/domain';

defineProps<{ visible: boolean; pending: boolean }>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  submit: [payload: SubmitResponsePayload];
}>();

const helpMethod = ref('');
const timeSlots = ref<string[]>([]);

const canSubmit = computed(() => helpMethod.value.trim().length > 0 && timeSlots.value.length > 0);

function submit() {
  if (!canSubmit.value) {
    return;
  }
  emit('submit', { helpMethod: helpMethod.value.trim(), timeSlots: [...timeSlots.value] });
}

function reset() {
  helpMethod.value = '';
  timeSlots.value = [];
}
</script>
