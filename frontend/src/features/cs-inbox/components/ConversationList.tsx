"use client";

import { MessageSquare } from "lucide-react";

import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Skeleton } from "@/shared/ui/Skeleton";

import { useConversations } from "../api/queries";

export function ConversationList() {
  const { data, isLoading, isError, error, refetch } = useConversations();

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 5 }).map((_, i) => (
          <Skeleton key={i} className="h-16 w-full rounded-lg" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <ErrorState
        message={error instanceof Error ? error.message : "Failed to load conversations"}
        onRetry={() => void refetch()}
      />
    );
  }

  if (!data?.data.length) {
    return (
      <EmptyState
        title="Inbox is clear"
        description="Customer conversations will appear here when messages arrive via WhatsApp."
        icon={<MessageSquare className="h-5 w-5" />}
      />
    );
  }

  return null;
}
