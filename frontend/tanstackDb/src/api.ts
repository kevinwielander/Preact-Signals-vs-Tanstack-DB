import type { Alarm, Resource, PaginatedResponse } from "./types";

const BASE = "http://localhost:8080";

export async function fetchMe(): Promise<Resource> {
  const res = await fetch(`${BASE}/me`);
  if (!res.ok) throw new Error("Failed to fetch /me");
  return res.json();
}

export async function fetchAlarms(
  page: number,
  pageSize: number,
  resourceId?: string
): Promise<PaginatedResponse<Alarm>> {
  const params = new URLSearchParams({
    page: String(page),
    pageSize: String(pageSize),
  });
  if (resourceId) params.set("resourceId", resourceId);
  const res = await fetch(`${BASE}/alarms?${params}`);
  if (!res.ok) throw new Error("Failed to fetch alarms");
  return res.json();
}

export async function fetchAlarm(id: string): Promise<Alarm> {
  const res = await fetch(`${BASE}/alarms/${id}`);
  if (!res.ok) throw new Error("Failed to fetch alarm");
  return res.json();
}

export async function fetchResources(): Promise<Resource[]> {
  const res = await fetch(`${BASE}/resources`);
  if (!res.ok) throw new Error("Failed to fetch resources");
  return res.json();
}

export async function patchAlarm(
  id: string,
  field: string,
  value: unknown,
  meId?: string
): Promise<Alarm> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (meId) headers["X-Resource-Id"] = meId;
  const res = await fetch(`${BASE}/alarms/${id}`, {
    method: "PATCH",
    headers,
    body: JSON.stringify({ [field]: value }),
  });
  if (!res.ok) throw new Error("Failed to patch alarm");
  return res.json();
}
