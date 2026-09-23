export interface Skill {
  id: number;
  owner: string;
  title: string;
  category: string;
  level: number;
  campus: string;
  description: string;
  timeSlots: string[];
  rewards: string[];
  portfolio: string;
}

export interface Need {
  id: number;
  requester: string;
  title: string;
  category: string;
  campus: string;
  expectTime: string;
  budgetType: string;
  description: string;
  responses: number;
  status: NeedStatus;
}

export type NeedStatus = 'open' | 'closed';

export type ResponseStatus = 'waiting' | 'interviewing' | 'withdrawn' | 'exited' | 'closed';

export interface NeedResponseEntry {
  id: number;
  needId: number;
  responder: string;
  helpMethod: string;
  timeSlots: string[];
  status: ResponseStatus;
  statusLabel: string;
  position: number;
  submittedAt: string;
  mine: boolean;
}

export interface ViewerQueueInfo {
  responseId: number;
  position: number;
  status: ResponseStatus;
  statusLabel: string;
  aheadCount: number;
  result: string;
}

export interface NeedSummary {
  id: number;
  requester: string;
  title: string;
  category: string;
  campus: string;
  expectTime: string;
  budgetType: string;
  description: string;
  responses: number;
  status: NeedStatus;
  statusLabel: string;
  responseCount: number;
  waitingCount: number;
  currentInterview: string;
  myQueue: ViewerQueueInfo | null;
}

export interface NeedDetail extends NeedSummary {
  queue: NeedResponseEntry[];
  canRespond: boolean;
  canPick: boolean;
  canClose: boolean;
}

export interface SubmitResponsePayload {
  helpMethod: string;
  timeSlots: string[];
}

export interface Match {
  id: number;
  provider: string;
  learner: string;
  offerSkill: string;
  wantedSkill: string;
  score: number;
  commonSlots: string[];
  recommendation: string;
}

export interface Appointment {
  id: number;
  pair: string;
  time: string;
  place: string;
  status: string;
  agenda: string;
}

export interface Review {
  id: number;
  from: string;
  to: string;
  rating: number;
  content: string;
}

export interface Conversation {
  id: number;
  withUser: string;
  unread: number;
  messages: string[];
}

export interface Profile {
  name: string;
  major: string;
  creditScore: number;
  creditLevel: string;
  skillWall: Skill[];
  radar: Record<string, number>;
  history: string[];
  reviews: Review[];
}

export interface Overview {
  service: string;
  categories: string[];
  metrics: Record<string, number>;
  skills: Skill[];
  needs: Need[];
  matches: Match[];
  appointments: Appointment[];
  reviews: Review[];
  messages: Conversation[];
  profile: Profile;
}
