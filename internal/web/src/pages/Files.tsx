import { Image } from "@kobalte/core/image";
import { createQuery } from "@tanstack/solid-query";
import { RiMediaImageLine, RiSystemDownloadLine } from "solid-icons/ri";
import { api } from "./data";
import { ErrorBoundary, For, Suspense } from "solid-js";
import { Skeleton } from "~/ui/Skeleton";
import { PageError, PageTitle } from "~/ui/Page";
import { formatDate, useQueryFilter, useQueryNumber } from "~/lib/utils";
import { LayoutCenter } from "~/ui/Layout";
import { DeviceFilterCombobox } from "~/components/Device";
import {
  PaginationEllipsis,
  PaginationEnd,
  PaginationItem,
  PaginationItems,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
  PaginationRoot,
  PaginationStart,
} from "~/ui/Pagination";

export default function Files() {
  const pageQuery = useQueryNumber("page", 1);
  const devicesFilter = useQueryFilter("devices");

  const data = createQuery(() => ({
    ...api.files.list({
      device: devicesFilter.values(),
    }),
    throwOnError: true,
  }));

  const devices = createQuery(() => ({
    ...api.devices.list,
    throwOnError: true,
  }));

  return (
    <LayoutCenter>
      <PageTitle>Files</PageTitle>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <div>
          <DeviceFilterCombobox
            deviceIDs={devicesFilter.values()}
            setDeviceIDs={devicesFilter.setValues}
          />
        </div>

        <PaginationRoot
          page={data.data?.pagination.page}
          count={data.data?.pagination.total_pages || 0}
          onPageChange={(page) => pageQuery.setValue(page)}
          itemComponent={(props) => (
            <PaginationItem page={props.page}>
              <PaginationLink
                isActive={props.page == data.data?.pagination.page}
              >
                {props.page}
              </PaginationLink>
            </PaginationItem>
          )}
          ellipsisComponent={() => <PaginationEllipsis />}
        >
          <PaginationStart>
            <PaginationItems />
          </PaginationStart>
          <PaginationEnd>
            <PaginationPrevious />
            <PaginationNext />
          </PaginationEnd>
        </PaginationRoot>

        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="grid grid-cols-1 gap-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
            <For each={devices.data}>
              {(device) => (
                <div class="w-full">
                  <div class="rounded-b border">
                    <Image>
                      <Image.Img
                        src={`/api/devices/${device.uuid}/snapshot`}
                        class="aspect-video object-contain"
                      />
                      <Image.Fallback>
                        <RiMediaImageLine class="aspect-video h-full w-full" />
                      </Image.Fallback>
                    </Image>
                    <div class="flex items-center justify-between gap-2 border-t p-2">
                      <div class="flex flex-col text-sm">
                        {formatDate(new Date())}
                      </div>
                      <RiSystemDownloadLine class="h-5 w-5" />
                    </div>
                  </div>
                </div>
              )}
            </For>
          </div>
        </Suspense>
      </ErrorBoundary>
    </LayoutCenter>
  );
}
