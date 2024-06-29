import { useQueryClient } from "@tanstack/solid-query";
import { api } from "./data";

export default function () {
  useQueryClient().prefetchQuery({
    ...api.settings.get,
  });
}
