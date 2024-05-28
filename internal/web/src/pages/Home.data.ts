import { useQueryClient } from "@tanstack/solid-query";
import { pages } from "./data";

export default function() {
  useQueryClient().prefetchQuery({
    ...pages.home
  })
}
