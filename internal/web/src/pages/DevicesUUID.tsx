import { A, useParams } from "@solidjs/router";
import { createMutation, createQuery } from "@tanstack/solid-query";
import {
  BreadcrumbsRoot,
  BreadcrumbsItem,
  BreadcrumbsLink,
  BreadcrumbsSeparator,
} from "~/ui/Breadcrumbs";
import { LayoutCenter } from "~/ui/Layout";
import { PageError, PageTitle } from "~/ui/Page";
import { api } from "./data";
import {
  ErrorBoundary,
  ParentProps,
  Show,
  Suspense,
  batch,
  createSignal,
} from "solid-js";
import { formatDate, parseDate } from "~/lib/utils";
import { Uptime } from "~/components/Utils";
import { Skeleton } from "~/ui/Skeleton";
import {
  TooltipArrow,
  TooltipContent,
  TooltipRoot,
  TooltipTrigger,
} from "~/ui/Tooltip";
import { createTimeAgo } from "@solid-primitives/date";
import { Image } from "@kobalte/core/image";
import { RiMediaImageLine } from "solid-icons/ri";
import { ToggleButton } from "@kobalte/core/toggle-button";
import { linkVariants } from "~/ui/Link";
import { Button } from "~/ui/Button";
import { postApiDevicesByUuidReboot } from "~/client";
import { toast } from "~/ui/Toast";
import {
  AlertDialogRoot,
  AlertDialogModal,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
} from "~/ui/AlertDialog";

export default function DevicesUUID() {
  const params = useParams<{ uuid: string }>();

  const data = createQuery(() => ({
    ...api.devices.get(params.uuid),
    throwOnError: true,
  }));
  const uptimeData = createQuery(() => ({
    ...api.devices.uptime(params.uuid),
    throwOnError: true,
  }));
  const detailData = createQuery(() => ({
    ...api.devices.detail(params.uuid),
    throwOnError: true,
  }));
  const softwareData = createQuery(() => ({
    ...api.devices.software(params.uuid),
    throwOnError: true,
  }));
  const licensesData = createQuery(() => ({
    ...api.devices.licenses(params.uuid),
    throwOnError: true,
  }));

  const [rebootDialog, setRebootDialog] = createSignal(false);
  const rebootMutation = createMutation(() => ({
    mutationFn: () => postApiDevicesByUuidReboot({ uuid: params.uuid }),
    onSuccess: () =>
      batch(() => {
        setRebootDialog(false);
        toast.show("Rebooting...");
      }),
    onError: (error) => toast.error(error.name, error.message),
  }));

  return (
    <LayoutCenter>
      <ErrorBoundary fallback={(error) => <PageError error={error} />}>
        <PageTitle>
          <BreadcrumbsRoot>
            <BreadcrumbsItem>
              <BreadcrumbsLink as={A} href="/devices">
                Devices
              </BreadcrumbsLink>
              <BreadcrumbsSeparator />
            </BreadcrumbsItem>
            <BreadcrumbsItem>
              <Suspense fallback={params.uuid}>{data.data?.name}</Suspense>
            </BreadcrumbsItem>
          </BreadcrumbsRoot>
        </PageTitle>

        <div class="flex justify-center gap-4">
          <div class="flex max-w-4xl flex-grow flex-col gap-4">
            <div class="rounded border">
              <Image>
                <Image.Img
                  src={`/api/devices/${params.uuid}/snapshot`}
                  alt=""
                />
                <Image.Fallback>
                  <RiMediaImageLine class="aspect-video h-full w-full" />
                </Image.Fallback>
              </Image>
            </div>

            <div class="rounded border p-2">
              <AlertDialogRoot
                open={rebootDialog()}
                onOpenChange={setRebootDialog}
              >
                <AlertDialogModal>
                  <AlertDialogHeader>
                    <AlertDialogTitle>
                      Are you sure you wish to reboot the device?
                    </AlertDialogTitle>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                      disabled={rebootMutation.isPending}
                      onClick={() => rebootMutation.mutate()}
                      variant="destructive"
                    >
                      Reboot
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogModal>
              </AlertDialogRoot>
              <Button
                disabled={rebootMutation.isPending}
                onClick={[setRebootDialog, true]}
              >
                Reboot
              </Button>
            </div>
          </div>

          <div class="max-w-lg flex-1">
            <div class="rounded border p-4">
              <PropertyTable>
                <PropertyRow name="Name">{data.data?.name}</PropertyRow>
                <PropertyRow name="IP">
                  <a href={"http://" + data.data?.ip} target="_blank">
                    {data.data?.ip}
                  </a>
                </PropertyRow>
                <PropertyRow name="Created At">
                  {formatDate(parseDate(data.data?.created_at))}
                </PropertyRow>
                <PropertyRow name="Updated At">
                  {formatDate(parseDate(data.data?.updated_at))}
                </PropertyRow>
                <Show when={data.data?.email}>
                  {(email) => <PropertyRow name="Email">{email()}</PropertyRow>}
                </Show>
                <Show when={data.data?.latitude}>
                  {(latitude) => (
                    <PropertyRow name="Latitude">{latitude()}</PropertyRow>
                  )}
                </Show>
                <Show when={data.data?.longitude}>
                  {(longitude) => (
                    <PropertyRow name="Longitude">{longitude()}</PropertyRow>
                  )}
                </Show>

                <ErrorBoundary
                  fallback={(error) => <PageError error={error} />}
                >
                  <Suspense
                    fallback={
                      <PropertyBlock>
                        <Skeleton class="h-6 w-full" />
                      </PropertyBlock>
                    }
                  >
                    <Show when={uptimeData.data?.supported}>
                      <PropertyRow name="Uptime">
                        <Uptime date={parseDate(uptimeData.data?.last)} />
                      </PropertyRow>
                    </Show>
                  </Suspense>

                  <Suspense
                    fallback={
                      <PropertyBlock>
                        <Skeleton class="h-32 w-full" />
                      </PropertyBlock>
                    }
                  >
                    <PropertyRow name="SN">
                      <ToggleButton>
                        {(state) => (
                          <Show
                            when={state.pressed()}
                            fallback={<>***************</>}
                          >
                            {detailData.data?.sn}
                          </Show>
                        )}
                      </ToggleButton>
                    </PropertyRow>
                    <PropertyRow name="Device Class">
                      {detailData.data?.device_class}
                    </PropertyRow>
                    <PropertyRow name="Device Type">
                      {detailData.data?.device_type}
                    </PropertyRow>
                    <PropertyRow name="Hardware Version">
                      {detailData.data?.hardware_version}
                    </PropertyRow>
                    <PropertyRow name="Market Area">
                      {detailData.data?.market_area}
                    </PropertyRow>
                    <PropertyRow name="Vendor">
                      {detailData.data?.vendor}
                    </PropertyRow>
                    <PropertyRow name="Onvif Version">
                      {detailData.data?.onvif_version}
                    </PropertyRow>
                  </Suspense>

                  <Suspense
                    fallback={
                      <PropertyBlock>
                        <Skeleton class="h-32 w-full" />
                      </PropertyBlock>
                    }
                  >
                    <PropertyRow name="Build">
                      {softwareData.data?.build}
                    </PropertyRow>
                    <PropertyRow name="Build Date">
                      {softwareData.data?.build_date}
                    </PropertyRow>
                    <PropertyRow name="Security Base Line Version">
                      {softwareData.data?.security_base_line_version}
                    </PropertyRow>
                    <PropertyRow name="Version">
                      {softwareData.data?.version}
                    </PropertyRow>
                    <PropertyRow name="Web Version">
                      {softwareData.data?.web_version}
                    </PropertyRow>
                  </Suspense>

                  <Suspense
                    fallback={
                      <PropertyBlock>
                        <Skeleton class="h-6 w-full" />
                      </PropertyBlock>
                    }
                  >
                    <Show when={licensesData.data?.at(0)}>
                      {(license) => {
                        const effectiveTime = () =>
                          parseDate(license().effective_time);
                        const [effectiveTimeAgo] = createTimeAgo(
                          effectiveTime,
                          {
                            interval: 0,
                          },
                        );

                        return (
                          <PropertyRow name={`License Effective Date`}>
                            <TooltipRoot>
                              <TooltipTrigger>
                                {formatDate(effectiveTime())}
                              </TooltipTrigger>
                              <TooltipContent>
                                <TooltipArrow />
                                {effectiveTimeAgo()}
                              </TooltipContent>
                            </TooltipRoot>
                          </PropertyRow>
                        );
                      }}
                    </Show>
                  </Suspense>
                </ErrorBoundary>
              </PropertyTable>
            </div>
          </div>
        </div>
      </ErrorBoundary>
    </LayoutCenter>
  );
}

function PropertyTable(props: ParentProps) {
  return <div class="flex flex-col gap-2 overflow-auto">{props.children}</div>;
}

function PropertyRow(props: ParentProps<{ name: string }>) {
  return (
    <div class="flex justify-between gap-2">
      <div>{props.name}</div>
      <div class="text-end text-muted-foreground">{props.children}</div>
    </div>
  );
}

function PropertyBlock(props: ParentProps) {
  return <div>{props.children}</div>;
}
