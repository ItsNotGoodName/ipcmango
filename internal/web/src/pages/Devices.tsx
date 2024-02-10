import Humanize from "humanize-plus"
import { A, createAsync, useSearchParams } from "@solidjs/router"
import { ErrorBoundary, For, Show, Suspense } from "solid-js"
import { PageError } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from "~/ui/Tabs"
import { TableBody, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table"
import { GetDevicesPageResp_Device } from "~/twirp/rpc"
import { getDeviceDetail, getDeviceRPCStatus, getDeviceSoftwareVersion, getListDeviceLicenses, getListDeviceStorage, } from "./data"
import { Skeleton } from "~/ui/Skeleton"
import { ToggleButton } from "@kobalte/core"
import { formatDate, parseDate } from "~/lib/utils"
import { getDevicesPage } from "./Devices.data"
import { linkVariants } from "~/ui/Link"

export function Devices() {
  const [searchParams, setSearchParams] = useSearchParams()
  const data = createAsync(() => getDevicesPage())

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <TabsRoot value={searchParams.tab || "device"} onChange={(value) => setSearchParams({ tab: value })}>
          <div class="overflow-x-auto">
            <TabsList>
              <TabsTrigger value="device" >Device</TabsTrigger>
              <TabsTrigger value="rpc-status" >RPC Status</TabsTrigger>
              <TabsTrigger value="detail" >Detail</TabsTrigger>
              <TabsTrigger value="software-version" >Software Version</TabsTrigger>
              <TabsTrigger value="license" >License</TabsTrigger>
              <TabsTrigger value="storage" >Storage</TabsTrigger>
            </TabsList>
          </div>
          <TabsContent value="device">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <DeviceTable devices={data()?.devices} />
            </Suspense>
          </TabsContent>
          <TabsContent value="rpc-status">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <RPCStatusTable devices={data()?.devices} />
            </Suspense>
          </TabsContent>
          <TabsContent value="detail">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <DetailTable devices={data()?.devices} />
            </Suspense>
          </TabsContent>
          <TabsContent value="software-version">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <SoftwareVersionTable devices={data()?.devices} />
            </Suspense>
          </TabsContent>
          <TabsContent value="license">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <LicenseTable devices={data()?.devices} />
            </Suspense>
          </TabsContent>
          <TabsContent value="storage">
            <Suspense fallback={<Skeleton class="h-32" />}>
              <StorageTable devices={data()?.devices} />
            </Suspense>
          </TabsContent>
        </TabsRoot>
      </ErrorBoundary>
    </LayoutNormal >
  )
}

function DeviceTable(props: { devices?: GetDevicesPageResp_Device[] }) {
  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>URL</TableHead>
          <TableHead>Username</TableHead>
          <TableHead>Disabled</TableHead>
          <TableHead>Created At</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => (
            <TableRow>
              <DeviceNameCell device={item} />
              <TableCell>
                <a class={linkVariants()} href={item.url}>{item.url}</a>
              </TableCell>
              <TableCell>
                {item.username}
              </TableCell>
              <TableCell>
                {item.disabled ? "TRUE" : "FALSE"}
              </TableCell>
              <TableCell>
                {formatDate(parseDate(item.createdAtTime))}
              </TableCell>
            </TableRow>
          )}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function RPCStatusTable(props: { devices?: GetDevicesPageResp_Device[] }) {
  const colspan = 8

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>State</TableHead>
          <TableHead>Last Login</TableHead>
          <TableHead>Error</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createAsync(() => getDeviceRPCStatus(item.id))

            return (
              <TableRow>
                <DeviceNameCell device={item} />
                <ErrorBoundary fallback={e => <ErrorTableCell colspan={colspan} error={e} />}>
                  <Suspense fallback={<LoadingTableCell colspan={colspan} />}>
                    <TableCell>{data()?.state}</TableCell>
                    <TableCell>{formatDate(parseDate(data()?.lastLoginTime))}</TableCell>
                    <TableCell>{data()?.error}</TableCell>
                  </Suspense>
                </ErrorBoundary>
              </TableRow>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function DetailTable(props: { devices?: GetDevicesPageResp_Device[] }) {
  const colspan = 9

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>SN</TableHead>
          <TableHead>Device Class</TableHead>
          <TableHead>Device Type</TableHead>
          <TableHead>Hardware Version</TableHead>
          <TableHead>Market Area</TableHead>
          <TableHead>Process Info</TableHead>
          <TableHead>Vendor</TableHead>
          <TableHead>Onvif Version</TableHead>
          <TableHead>Algorithm Version</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createAsync(() => getDeviceDetail(item.id))

            return (
              <TableRow>
                <DeviceNameCell device={item} />
                <ErrorBoundary fallback={e => <ErrorTableCell colspan={colspan} error={e} />}>
                  <Suspense fallback={<LoadingTableCell colspan={colspan} />}>
                    <TableCell>
                      <ToggleButton.Root>
                        {state => (
                          <Show when={state.pressed()} fallback={<>***************</>}>
                            {data()?.sn}
                          </Show>
                        )}
                      </ToggleButton.Root>
                    </TableCell>
                    <TableCell>{data()?.deviceClass}</TableCell>
                    <TableCell>{data()?.deviceType}</TableCell>
                    <TableCell>{data()?.hardwareVersion}</TableCell>
                    <TableCell>{data()?.marketArea}</TableCell>
                    <TableCell>{data()?.processInfo}</TableCell>
                    <TableCell>{data()?.vendor}</TableCell>
                    <TableCell>{data()?.onvifVersion}</TableCell>
                    <TableCell>{data()?.algorithmVersion}</TableCell>
                  </Suspense>
                </ErrorBoundary>
              </TableRow>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function SoftwareVersionTable(props: { devices?: GetDevicesPageResp_Device[] }) {
  const colspan = 9

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Build</TableHead>
          <TableHead>Build Date</TableHead>
          <TableHead>Security Base Line Version</TableHead>
          <TableHead>Version</TableHead>
          <TableHead>Web Version</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createAsync(() => getDeviceSoftwareVersion(item.id))

            return (
              <TableRow>
                <DeviceNameCell device={item} />
                <ErrorBoundary fallback={e => <ErrorTableCell colspan={colspan} error={e} />}>
                  <Suspense fallback={<LoadingTableCell colspan={colspan} />}>
                    <TableCell>{data()?.build}</TableCell>
                    <TableCell>{data()?.buildDate}</TableCell>
                    <TableCell>{data()?.securityBaseLineVersion}</TableCell>
                    <TableCell>{data()?.version}</TableCell>
                    <TableCell>{data()?.webVersion}</TableCell>
                  </Suspense>
                </ErrorBoundary>
              </TableRow>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function LicenseTable(props: { devices?: GetDevicesPageResp_Device[] }) {
  const colspan = 9

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Abroad Info</TableHead>
          <TableHead>All Type</TableHead>
          <TableHead>Digit Channel</TableHead>
          <TableHead>Effective Days</TableHead>
          <TableHead>Effective Time</TableHead>
          <TableHead>License ID</TableHead>
          <TableHead>Product Type</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Username</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createAsync(() => getListDeviceLicenses(item.id))

            return (
              <ErrorBoundary fallback={e =>
                <TableRow>
                  <DeviceNameCell device={item} />
                  <ErrorTableCell colspan={colspan} error={e} />
                </TableRow>
              }>
                <Suspense fallback={
                  <TableRow>
                    <DeviceNameCell device={item} />
                    <LoadingTableCell colspan={colspan} />
                  </TableRow>
                }>
                  <For each={data()} fallback={
                    <TableRow>
                      <DeviceNameCell device={item} />
                      <TableCell colspan={colspan}>N/A</TableCell>
                    </TableRow>
                  }>
                    {v => (
                      <TableRow>
                        <DeviceNameCell device={item} />
                        <TableCell>{v.abroadInfo}</TableCell>
                        <TableCell>{v.allType}</TableCell>
                        <TableCell>{v.digitChannel}</TableCell>
                        <TableCell>{v.effectiveDays}</TableCell>
                        <TableCell>{formatDate(parseDate(v.effectiveTime))}</TableCell>
                        <TableCell>{v.licenseId}</TableCell>
                        <TableCell>{v.productType}</TableCell>
                        <TableCell>{v.status}</TableCell>
                        <TableCell>{v.username}</TableCell>
                      </TableRow>
                    )
                    }
                  </For>
                </Suspense>
              </ErrorBoundary>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function StorageTable(props: { devices?: GetDevicesPageResp_Device[] }) {
  const colspan = 6;

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Name</TableHead>
          <TableHead>State</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Used</TableHead>
          <TableHead>Total</TableHead>
          <TableHead>Is Error</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createAsync(() => getListDeviceStorage(item.id))

            return (
              <ErrorBoundary fallback={e =>
                <TableRow>
                  <DeviceNameCell device={item} />
                  <ErrorTableCell colspan={colspan} error={e} />
                </TableRow>
              }>
                <Suspense fallback={
                  <TableRow>
                    <DeviceNameCell device={item} />
                    <LoadingTableCell colspan={colspan} />
                  </TableRow>
                }>
                  <For each={data()} fallback={
                    <TableRow>
                      <DeviceNameCell device={item} />
                      <TableCell colspan={colspan}>N/A</TableCell>
                    </TableRow>
                  }>
                    {v => (
                      <TableRow>
                        <DeviceNameCell device={item} />
                        <TableCell>{v.name}</TableCell>
                        <TableCell>{v.state}</TableCell>
                        <TableCell>{v.type}</TableCell>
                        <TableCell>{Humanize.fileSize(Number(v.usedBytes))}</TableCell>
                        <TableCell>{Humanize.fileSize(Number(v.totalBytes))}</TableCell>
                        <TableCell>{v.isError}</TableCell>
                      </TableRow>
                    )
                    }
                  </For>
                </Suspense>
              </ErrorBoundary>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function LoadingTableCell(props: { colspan: number }) {
  return (
    <TableCell colspan={props.colspan} class="py-0">
      <Skeleton class="h-8" />
    </TableCell>
  )
}

function ErrorTableCell(props: { colspan: number, error: Error }) {
  return (
    <TableCell colspan={props.colspan} class="py-0">
      <div class="bg-destructive text-destructive-foreground rounded p-2">
        {props.error.message}
      </div>
    </TableCell>
  )
}

function DeviceNameCell(props: { device: { id: bigint, name: string } }) {
  return (
    <TableCell>
      <A class={linkVariants()} href={`./${props.device.id}`}>{props.device.name}</A>
    </TableCell>
  )
}
