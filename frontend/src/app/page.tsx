import { redirect } from "next/navigation";

import { ROUTES } from "@/shared/config/constants";

export default function HomePage() {
  redirect(ROUTES.login);
}
