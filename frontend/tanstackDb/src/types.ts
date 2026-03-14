export interface Alarm {
  id: string;
  title: string;
  description: string;
  severity: "low" | "medium" | "high" | "critical";
  status: "active" | "acknowledged" | "resolved";
  assignedResources: string[];
  eventHash: string;
  eventNumber: number;
  archivedOnOffset: number | null;
  createdOnOffset: number;
  createdAt: string;
  updatedAt: string;
}

export interface Resource {
  id: string;
  displayName: string;
  email: string;
  isUserAssociated: boolean;
  thumbnail: string;
  eventHash: string;
  eventNumber: number;
  archivedOnOffset: number | null;
  createdOnOffset: number;
  createdAt: string;
  updatedAt: string;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}
