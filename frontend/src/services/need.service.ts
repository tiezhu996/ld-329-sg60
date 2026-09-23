import type { NeedDetail, NeedSummary, SubmitResponsePayload } from '../types/domain';

const API_BASE = '/api';

interface ErrorBody {
  code?: string;
  message?: string;
}

async function parseError(response: Response): Promise<never> {
  let message = '请求失败，请稍后重试';
  try {
    const body = (await response.json()) as ErrorBody;
    if (body.message) {
      message = body.message;
    }
  } catch {
    // 保留默认错误文案
  }
  throw new Error(message);
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, init);
  if (!response.ok) {
    await parseError(response);
  }
  return response.json() as Promise<T>;
}

function post<T>(path: string, body: unknown): Promise<T> {
  return request<T>(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

export function fetchNeeds(user: string): Promise<NeedSummary[]> {
  return request<NeedSummary[]>(`/needs?user=${encodeURIComponent(user)}`);
}

export function fetchNeedDetail(needId: number, user: string): Promise<NeedDetail> {
  return request<NeedDetail>(`/needs/${needId}?user=${encodeURIComponent(user)}`);
}

export function submitNeedResponse(needId: number, operator: string, payload: SubmitResponsePayload): Promise<NeedDetail> {
  return post<NeedDetail>(`/needs/${needId}/responses`, { operator, ...payload });
}

export function pickNeedCandidate(needId: number, operator: string, responseId: number): Promise<NeedDetail> {
  return post<NeedDetail>(`/needs/${needId}/pick`, { operator, responseId });
}

export function exitNeedResponse(needId: number, operator: string): Promise<NeedDetail> {
  return post<NeedDetail>(`/needs/${needId}/exit`, { operator });
}

export function closeNeed(needId: number, operator: string): Promise<NeedDetail> {
  return post<NeedDetail>(`/needs/${needId}/close`, { operator });
}
