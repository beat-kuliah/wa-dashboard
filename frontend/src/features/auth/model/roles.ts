import type { User, UserRole } from "@/entities/user";

export function userHasRole(user: User | null | undefined, role: UserRole): boolean {
  return user?.roles.includes(role) ?? false;
}

export function isTenantAdmin(user: User | null | undefined): boolean {
  return userHasRole(user, "admin");
}

export function canViewTenantMembers(user: User | null | undefined): boolean {
  return userHasRole(user, "admin") || userHasRole(user, "supervisor");
}
