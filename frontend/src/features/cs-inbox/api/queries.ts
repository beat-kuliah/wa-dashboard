import { useQuery } from "@tanstack/react-query";

import { fetchConversations } from "./inboxApi";
import type { ConversationListParams } from "../model/schema";

export const inboxKeys = {
  all: ["inbox"] as const,
  conversations: (params: ConversationListParams) =>
    [...inboxKeys.all, "conversations", params] as const,
};

export function useConversations(params: ConversationListParams = {}) {
  return useQuery({
    queryKey: inboxKeys.conversations(params),
    queryFn: () => fetchConversations(params),
  });
}
