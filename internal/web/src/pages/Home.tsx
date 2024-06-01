import Humanize from "humanize-plus"
import { A, } from "@solidjs/router"
import { CardRoot, } from "~/ui/Card"
import { ErrorBoundary, ParentProps, Suspense } from "solid-js"
import { BiRegularCctv } from "solid-icons/bi"
import { PageError, PageLoading } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { RiBusinessMailLine, RiDeviceDatabase2Line, RiDeviceHardDrive2Line, RiDocumentFile2Line, RiWeatherFlashlightLine } from "solid-icons/ri"
import { formatDate, parseDate } from "~/lib/utils"
import { linkVariants } from "~/ui/Link"
import { createQuery } from "@tanstack/solid-query"
import { pages } from "./data"

function StatParent(props: ParentProps) {
  return <div class="flex-1 sm:max-w-48">{props.children}</div>
}

function StatRoot(props: ParentProps) {
  return <CardRoot class="flex gap-2 p-4">{props.children}</CardRoot>
}

function StatTitle(props: ParentProps) {
  return <h2>{props.children}</h2>
}

function StatValue(props: ParentProps) {
  return <p class="text-lg font-bold">{props.children}</p>
}

export default function () {
  const data = createQuery(() => ({
    ...pages.home,
    throwOnError: true
  }))

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<PageLoading />}>
          <div class="flex flex-col flex-wrap gap-4 sm:flex-row">
            <StatParent>
              <StatRoot>
                <A href="/devices" class="flex items-center">
                  <BiRegularCctv class="w-8 h-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Devices</StatTitle>
                  <StatValue>{data.data?.device_count}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <A href="/emails" class="flex items-center">
                  <RiBusinessMailLine class="w-8 h-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Emails</StatTitle>
                  <StatValue>{data.data?.email_count}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <A href="/events" class="flex items-center">
                  <RiWeatherFlashlightLine class="w-8 h-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Events</StatTitle>
                  <StatValue>{data.data?.event_count}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <A href="/files" class="flex items-center">
                  <RiDocumentFile2Line class="w-8 h-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Files</StatTitle>
                  <StatValue>{data.data?.file_count}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <div class="flex items-center">
                  <RiDeviceHardDrive2Line class="w-8 h-8" />
                </div>
                <div class="flex-1">
                  <StatTitle>File usage</StatTitle>
                  <StatValue>{Humanize.fileSize(data.data?.file_usage || 0)}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <div class="flex items-center">
                  <RiDeviceDatabase2Line class="w-8 h-8" />
                </div>
                <div class="flex-1">
                  <StatTitle>DB usage</StatTitle>
                  <StatValue>{Humanize.fileSize(data.data?.db_usage || 0)}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
          </div>
          {/* <div class="flex flex-col gap-4 lg:flex-row"> */}
          {/*   <div class="flex-1 lg:max-w-sm"> */}
          {/*     <CardRoot class="p-4"> */}
          {/*       <h1>Latest emails</h1> */}
          {/*       <div> */}
          {/*         <For each={data.data?.emails}> */}
          {/*           {v => { */}
          {/*             const [createdAt] = createDate(() => parseDate(v.createdAtTime)); */}
          {/*             const [createdAtAgo] = createTimeAgo(createdAt); */}
          {/**/}
          {/*             return ( */}
          {/*               <div class="flex flex-col border-b transition-colors sm:flex-row hover:bg-muted/50"> */}
          {/*                 <A href={`/emails/${v.id}`} class="flex flex-col flex-1 gap-2 p-2 sm:flex-row sm:pr-1 max-sm:pb-1"> */}
          {/*                   <div class="flex sm:min-w-32"> */}
          {/*                     <TooltipRoot> */}
          {/*                       <TooltipTrigger class="text-sm font-bold truncate text-start">{createdAtAgo()}</TooltipTrigger> */}
          {/*                       <TooltipContent> */}
          {/*                         <TooltipArrow /> */}
          {/*                         {formatDate(createdAt())} */}
          {/*                       </TooltipContent> */}
          {/*                     </TooltipRoot> */}
          {/*                   </div> */}
          {/*                   <div class="flex-1 truncate"> */}
          {/*                     {v.subject} */}
          {/*                   </div> */}
          {/*                 </A> */}
          {/*                 <Show when={v.attachmentCount > 0}> */}
          {/*                   <A href={`/emails/${v.id}?tab=attachments`} class="p-2 sm:pl-1 max-sm:pt-1"> */}
          {/*                     <TooltipRoot> */}
          {/*                       <TooltipTrigger class="flex items-center h-full"> */}
          {/*                         <RiEditorAttachment2 class="w-5 h-5" /> */}
          {/*                       </TooltipTrigger> */}
          {/*                       <TooltipContent> */}
          {/*                         <TooltipArrow /> */}
          {/*                         {v.attachmentCount} {Humanize.pluralize(v.attachmentCount, "attachment")} */}
          {/*                       </TooltipContent> */}
          {/*                     </TooltipRoot> */}
          {/*                   </A> */}
          {/*                 </Show> */}
          {/*               </div> */}
          {/*             ) */}
          {/*           }} */}
          {/*         </For> */}
          {/*       </div> */}
          {/*     </CardRoot> */}
          {/*   </div> */}
          {/*   <div class="flex flex-col flex-1 gap-4"> */}
          {/*     <h1>Latest files</h1> */}
          {/*     <div class="grid grid-cols-2 gap-4 sm:grid-cols-4 xl:grid-cols-6 2xl:grid-cols-8"> */}
          {/*       <For each={[]}> */}
          {/*         {(v) => { */}
          {/*           const [startTime] = createDate(() => parseDate(v.startTime)); */}
          {/*           const [startTimeAgo] = createTimeAgo(startTime); */}
          {/**/}
          {/*           return ( */}
          {/*             <div> */}
          {/*               <div class="flex flex-col w-full rounded-b border transition-all sm:max-w-48 hover:bg-accent/50"> */}
          {/*                 <A href={`/files/${v.id}`} > */}
          {/*                   <Image.Root class="mx-auto w-full max-h-48"> */}
          {/*                     <Image.Img src={v.thumbnailUrl} class="object-contain w-full h-full" /> */}
          {/*                     <Image.Fallback> */}
          {/*                       <Show when={v.type == "jpg"} fallback={ */}
          {/*                         <RiMediaVideoLine class="object-contain w-full h-full" /> */}
          {/*                       }> */}
          {/*                         <RiMediaImageLine class="object-contain w-full h-full" /> */}
          {/*                       </Show> */}
          {/*                     </Image.Fallback> */}
          {/*                   </Image.Root> */}
          {/*                 </A> */}
          {/*                 <Seperator /> */}
          {/*                 <div class="flex gap-2 justify-between items-center p-2"> */}
          {/*                   <TooltipRoot> */}
          {/*                     <TooltipTrigger class="text-sm truncate">{startTimeAgo()}</TooltipTrigger> */}
          {/*                     <TooltipContent> */}
          {/*                       <TooltipArrow /> */}
          {/*                       {formatDate(startTime())} */}
          {/*                     </TooltipContent> */}
          {/*                   </TooltipRoot> */}
          {/*                   <a href={v.url} target="_blank" title="Download"> */}
          {/*                     <RiSystemDownloadLine class="w-5 h-5" /> */}
          {/*                   </a> */}
          {/*                 </div> */}
          {/*               </div> */}
          {/*             </div> */}
          {/*           ) */}
          {/*         }} */}
          {/*       </For> */}
          {/*     </div> */}
          {/*   </div> */}
          {/* </div> */}
          <div class="flex flex-col sm:flex-row">
            <CardRoot class="p-4">
              <h1>Build</h1>
              <div class="overflow-x-auto relative">
                <table class="w-full">
                  <tbody>
                    <tr class="border-b">
                      <td class="p-2">Commit</td>
                      <td class="p-2"><a class={linkVariants()} href={data.data?.build.commit_url}>{data.data?.build?.commit}</a></td>
                    </tr>
                    <tr class="border-b">
                      <td class="p-2">Date</td>
                      <td class="p-2">{formatDate(parseDate(data.data?.build.date))}</td>
                    </tr>
                    <tr class="border-b">
                      <td class="p-2">Version</td>
                      <td class="p-2"><a class={linkVariants()} href={data.data?.build.release_url}>{data.data?.build?.version}</a></td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </CardRoot>
          </div>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal >
  )
}
