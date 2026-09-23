import { defineStore } from 'pinia';
import {
  closeNeed,
  fetchNeedDetail,
  fetchNeeds,
  pickInterview,
  submitNeedResponse,
  withdrawResponse,
} from '../services/need.service';
import type { NeedDetail, NeedSummary, SubmitResponsePayload } from '../types/domain';

// 技能求助候选名单：列表、详情和所有排队/约谈操作的状态管理
export const useNeedsStore = defineStore('needs', {
  state: () => ({
    viewer: '',
    list: [] as NeedSummary[],
    detail: null as NeedDetail | null,
    listLoading: false,
    detailLoading: false,
    acting: false,
  }),
  actions: {
    setViewer(viewer: string) {
      this.viewer = viewer;
    },
    async loadList() {
      this.listLoading = true;
      try {
        this.list = await fetchNeeds(this.viewer);
      } finally {
        this.listLoading = false;
      }
    },
    async loadDetail(needId: number) {
      this.detailLoading = true;
      try {
        this.detail = await fetchNeedDetail(needId, this.viewer);
      } finally {
        this.detailLoading = false;
      }
    },
    // 每次操作后以后端返回的最新详情刷新，并同步列表中的轮候信息
    async applyMutation(mutation: Promise<NeedDetail>) {
      this.acting = true;
      try {
        this.detail = await mutation;
        await this.loadList();
      } finally {
        this.acting = false;
      }
    },
    async submit(payload: SubmitResponsePayload) {
      if (!this.detail) return;
      await this.applyMutation(submitNeedResponse(this.detail.id, payload));
    },
    async pick(responseId: number) {
      if (!this.detail) return;
      await this.applyMutation(pickInterview(this.detail.id, responseId, this.viewer));
    },
    async withdraw(responseId: number) {
      if (!this.detail) return;
      await this.applyMutation(withdrawResponse(this.detail.id, responseId, this.viewer));
    },
    async close() {
      if (!this.detail) return;
      await this.applyMutation(closeNeed(this.detail.id, this.viewer));
    },
  },
});
