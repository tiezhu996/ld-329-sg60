import { AppException } from '../errors/AppException';
import { logger } from '../logger/logger';
import type { NeedDetail, NeedSummary, SubmitResponsePayload } from '../types/domain';

const API_BASE = '/api';

interface ErrorBody {
  code?: string;
  message?: string;
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, init);
  if (!response.ok) {
    const body = (await response.json().catch(() => ({}))) as ErrorBody;
    logger.error('求助接口请求失败', path, body);
    throw new AppException(body.message ?? '操作失败，请稍后重试', body.code ?? 'NEED_API_ERROR');
  }
  return response.json() as Promise<T>;
}

function post<T>(path: string, payload: unknown): Promise<T> {
  return request<T>(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
}

export function fetchNeeds(viewer: string): Promise<NeedSummary[]> {
  return request<NeedSummary[]>(`/needs?viewer=${encodeURIComponent(viewer)}`);
}

export function fetchNeedDetail(needId: number, viewer: string): Promise<NeedDetail> {
  return request<NeedDetail>(`/needs/${needId}?viewer=${encodeURIComponent(viewer)}`);
}

export function submitNeedResponse(needId: number, payload: SubmitResponsePayload): Promise<NeedDetail> {
  return post<NeedDetail>(`/needs/${needId}/responses`, payload);
}

export function pickInterview(needId: number, responseId: number, operator: string): Promise<NeedDetail> {
  return post<NeedDetail>(`/needs/${needId}/responses/${responseId}/interview`, { operator });
}

export function withdrawResponse(needId: number, responseId: number, operator: string): Promise<NeedDetail> {
  return post<NeedDetail>(`/needs/${needId}/responses/${responseId}/withdraw`, { operator });
}

export function closeNeed(needId: number, operator: string): Promise<NeedDetail> {
  return post<NeedDetail>(`/needs/${needId}/close`, { operator });
}
