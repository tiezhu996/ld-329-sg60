import type { NeedStatus, ResponseStatus } from '../types/domain';

type TagType = 'primary' | 'success' | 'info' | 'warning' | 'danger';

export const NEED_STATUS_TAG: Record<NeedStatus, TagType> = {
  open: 'success',
  closed: 'info',
};

export const RESPONSE_STATUS_TAG: Record<ResponseStatus, TagType> = {
  waiting: 'info',
  interviewing: 'warning',
  withdrawn: 'danger',
  exited: 'info',
  closed: 'info',
};

// 可上门时段选项（响应表单勾选）
export const VISIT_TIME_SLOTS = [
  '周一晚',
  '周二晚',
  '周三晚',
  '周四晚',
  '周五晚',
  '周六上午',
  '周六下午',
  '周日上午',
  '周日下午',
] as const;

export const NEED_BOARD_TITLE = '技能求助 · 候选名单';
export const NO_INTERVIEW_TEXT = '暂无约谈';
export const NOT_JOINED_TEXT = '未参与';
