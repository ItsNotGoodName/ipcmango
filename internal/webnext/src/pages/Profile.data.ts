import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getProfile = cache(() => useClient().page.profile({}).then((req) => req.response), "profile")

export default function() {
  void getProfile()
}

