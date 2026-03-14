import {
  useQuery,
  useMutation,
  useQueryClient,
  keepPreviousData,
} from "@tanstack/react-query";
import { fetchMe, fetchAlarms, fetchAlarm, fetchResources, patchAlarm } from "./api";

export function useMe() {
  return useQuery({
    queryKey: ["me"],
    queryFn: fetchMe,
    staleTime: Infinity,
  });
}

export function useResources() {
  return useQuery({
    queryKey: ["resources"],
    queryFn: fetchResources,
    staleTime: Infinity,
  });
}

export function useResourceMap() {
  const { data: resources } = useResources();
  const map = new Map<string, string>();
  if (resources) {
    for (const r of resources) {
      map.set(r.id, r.displayName);
    }
  }
  return map;
}

export function useMyAlarms(meId: string | undefined) {
  return useQuery({
    queryKey: ["alarms", "mine", meId],
    queryFn: () => fetchAlarms(1, 1000, meId!),
    enabled: !!meId,
  });
}

export function useAlarmsPaginated(page: number, pageSize: number) {
  return useQuery({
    queryKey: ["alarms", "paginated", page, pageSize],
    queryFn: () => fetchAlarms(page, pageSize),
    placeholderData: keepPreviousData,
  });
}

export function useAlarm(id: string | null) {
  return useQuery({
    queryKey: ["alarm", id],
    queryFn: () => fetchAlarm(id!),
    enabled: !!id,
  });
}

export function usePatchAlarm() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      field,
      value,
    }: {
      id: string;
      field: string;
      value: unknown;
    }) => patchAlarm(id, field, value),
    onSuccess: (_data, { id }) => {
      queryClient.invalidateQueries({ queryKey: ["alarm", id] });
      queryClient.invalidateQueries({ queryKey: ["alarms"] });
    },
  });
}
