import type { MyQueueStatus, NeedStatus, ResponseStatus } from '../types/domain';

export const NEED_STATUS_LABELS: Record<NeedStatus, string> = {
  open: '开放中',
  closed: '已关闭',
};

export const NEED_STATUS_TAG_TYPES: Record<NeedStatus, 'success' | 'info'> = {
  open: 'success',
  closed: 'info',
};

export const RESPONSE_STATUS_LABELS: Record<ResponseStatus, string> = {
  waiting: '排队中',
  interviewing: '约谈中',
  withdrawn: '已退出',
};

export const RESPONSE_STATUS_TAG_TYPES: Record<ResponseStatus, 'warning' | 'success' | 'info'> = {
  waiting: 'warning',
  interviewing: 'success',
  withdrawn: 'info',
};

export const MY_STATUS_LABELS: Record<MyQueueStatus, string> = {
  none: '未响应',
  waiting: '排队中',
  interviewing: '约谈中',
  withdrawn: '已退出',
};

export const VISIT_SLOT_OPTIONS = [
  '周一晚',
  '周二晚',
  '周三晚',
  '周四晚',
  '周五晚',
  '周六上午',
  '周六下午',
  '周日全天',
] as const;

export const NEED_QUEUE_MESSAGES = {
  submitSuccess: '已加入候选名单，按提交先后排队',
  pickSuccess: '已安排约谈，其他响应者继续等待',
  withdrawSuccess: '已退出轮候，由下一位自动接手',
  closeSuccess: '求助已关闭，候选名单停止',
  loadFailed: '无法加载求助列表',
  detailFailed: '无法加载求助详情',
  needHelpOffer: '请填写能帮的方式',
  needVisitSlots: '请至少选择一个可上门时段',
} as const;
