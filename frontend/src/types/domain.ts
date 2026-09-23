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

export type NeedStatus = 'open' | 'closed';
export type ResponseStatus = 'waiting' | 'interviewing' | 'withdrawn';
export type MyQueueStatus = 'none' | ResponseStatus;

export interface Need {
  id: number;
  requester: string;
  title: string;
  category: string;
  campus: string;
  expectTime: string;
  budgetType: string;
  description: string;
  status: NeedStatus;
  responses: number;
}

export interface NeedSummary extends Need {
  currentInterviewee: string;
  waitingCount: number;
  isPublisher: boolean;
  myStatus: MyQueueStatus;
  myPosition: number;
}

export interface NeedQueueEntry {
  responseId: number;
  position: number;
  responder: string;
  helpOffer: string;
  visitSlots: string[];
  status: ResponseStatus;
  submittedAt: string;
  isViewer: boolean;
}

export interface NeedDetail extends NeedSummary {
  entries: NeedQueueEntry[];
}

export interface SubmitResponsePayload {
  responder: string;
  helpOffer: string;
  visitSlots: string[];
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
  needs: NeedSummary[];
  matches: Match[];
  appointments: Appointment[];
  reviews: Review[];
  messages: Conversation[];
  profile: Profile;
}
