import { apiClient, buildQueryString } from "@/shared/lib/apiClient";

import {
  conversationListResponseSchema,
  type ConversationListParams,
} from "../model/schema";

export async function fetchConversations(params: ConversationListParams = {}) {
  const qs = buildQueryString({
    limit: params.limit ?? 20,
    offset: params.offset ?? 0,
    status: params.status,
  });

  const json = await apiClient<unknown>(`/inbox/conversations${qs}`);
  return conversationListResponseSchema.parse(json);
}
