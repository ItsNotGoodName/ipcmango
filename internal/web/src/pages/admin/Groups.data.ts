import { cache } from "@solidjs/router"
import { parseOrder } from "~/lib/utils"
import { PageProps } from "~/lib/utils"
import { useClient } from "~/providers/client"
import { GetAdminGroupsPageReq } from "~/twirp/rpc"

export const getAdminGroupsPage = cache((input: GetAdminGroupsPageReq) => useClient().admin.getAdminGroupsPage(input).then((req) => req.response), "getAdminGroupsPage")

export type AdminGroupsPageSearchParams = {
  page: string
  perPage: string
  sort: string
  order: string
  filter: string
}

export default function({ params }: PageProps<AdminGroupsPageSearchParams>) {
  void getAdminGroupsPage({
    page: {
      page: Number(params.page) || 1,
      perPage: Number(params.perPage) || 10
    },
    sort: {
      field: params.sort || "",
      order: parseOrder(params.order),
    },
  })
}

export const getGroup = cache((id: bigint) => useClient().admin.getGroup({ id: id }).then((req) => req.response), "getGroup")
