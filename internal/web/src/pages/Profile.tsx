import Humanize from "humanize-plus"
import { action, createAsync, revalidate, useAction, useSubmission } from "@solidjs/router"
import { RiSystemCheckLine, RiSystemCloseLine } from "solid-icons/ri"
import { ErrorBoundary, For, ParentProps, Show, Suspense, createSignal, } from "solid-js"
import { createForm, required, reset } from "@modular-forms/solid"

import { formatDate, parseDate, catchAsToast, throwAsFormError } from "~/lib/utils"
import { CardRoot, } from "~/ui/Card"
import { getProfilePage } from "./Profile.data"
import { Button } from "~/ui/Button"
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table"
import { useClient } from "~/providers/client"
import { Badge } from "~/ui/Badge"
import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form"
import { Input } from "~/ui/Input"
import { Skeleton } from "~/ui/Skeleton"
import { getSession } from "~/providers/session"
import { PageError } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { AlertDialogAction, AlertDialogCancel, AlertDialogModal, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle } from "~/ui/AlertDialog"
import { Shared } from "~/components/Shared"

const actionRevokeAllMySessions = action(() => useClient()
  .user.revokeAllMySessions({})
  .then(() => revalidate(getProfilePage.key))
  .catch(catchAsToast))

const actionRevokeMySession = action((sessionId: bigint) => useClient()
  .user.revokeMySession({ sessionId })
  .then(() => revalidate(getProfilePage.key))
  .catch(catchAsToast))

export function Profile() {
  const data = createAsync(() => getProfilePage())

  const [revokeAllMySessionsConfirm, setRevokeAllMySessionsConfirm] = createSignal(false)
  const revokeAllMySessionsSubmission = useSubmission(actionRevokeAllMySessions)
  const revokeAllMySessionsAction = useAction(actionRevokeAllMySessions)
  const revokeAllMySessions = () => revokeAllMySessionsAction()
    .then(() => setRevokeAllMySessionsConfirm(false))

  const [revokeMySessionsConfirm, setRevokeMySessionsConfirm] = createSignal(BigInt(0))
  const revokeMySessionSubmission = useSubmission(actionRevokeMySession)
  const revokeMySessionAction = useAction(actionRevokeMySession)
  const revokeMySession = () => revokeMySessionAction(revokeMySessionsConfirm()).
    then(() => setRevokeMySessionsConfirm(BigInt(0)))

  return (
    <LayoutNormal class="max-w-4xl">
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>

        <AlertDialogRoot open={revokeAllMySessionsConfirm()} onOpenChange={setRevokeAllMySessionsConfirm}>
          <AlertDialogModal>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you sure you wish to revoke all sessions?</AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction disabled={revokeAllMySessionsSubmission.pending} onClick={revokeAllMySessions} variant="destructive">
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogModal>
        </AlertDialogRoot>

        <AlertDialogRoot open={revokeMySessionsConfirm() != BigInt(0)} onOpenChange={() => setRevokeMySessionsConfirm(BigInt(0))}>
          <AlertDialogModal>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you sure you wish to revoke this session?</AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction disabled={revokeMySessionSubmission.pending} onClick={revokeMySession} variant="destructive">
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogModal>
        </AlertDialogRoot>

        <Shared.Title>Profile</Shared.Title>

        <CardRoot class="overflow-x-auto p-4">
          <Suspense fallback={<Skeleton class="h-32" />}>
            <table>
              <tbody>
                <tr>
                  <th class="pr-2">Username</th>
                  <td>{data()?.username}</td>
                </tr>
                <tr>
                  <th class="pr-2">Email</th>
                  <td>{data()?.email}</td>
                </tr>
                <tr>
                  <th class="pr-2">Admin</th>
                  <td>
                    <Show when={data()?.admin} fallback={<RiSystemCloseLine class="h-6 w-6 text-red-500" />}>
                      <RiSystemCheckLine class="h-6 w-6 text-green-500" />
                    </Show>
                  </td>
                </tr>
                <tr>
                  <th class="pr-2">Created At</th>
                  <td>{formatDate(parseDate(data()?.createdAtTime))}</td>
                </tr>
                <tr>
                  <th class="pr-2">Updated At</th>
                  <td>{formatDate(parseDate(data()?.updatedAtTime))}</td>
                </tr>
              </tbody>
            </table>
          </Suspense>
        </CardRoot>

        <Shared.Title>Change username</Shared.Title>
        <Center>
          <ChangeUsernameForm />
        </Center>

        <Shared.Title>Change password</Shared.Title>
        <Center>
          <ChangePasswordForm />
        </Center>

        <Shared.Title>Sessions</Shared.Title>
        <div class="flex">
          <Button variant="destructive" onClick={() => setRevokeAllMySessionsConfirm(true)}>
            Revoke all sessions
          </Button>
        </div>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <TableRoot>
            <TableCaption>{data()?.sessions.length} {Humanize.pluralize(data()?.sessions.length || 0, "Session")}</TableCaption>
            <TableHeader>
              <TableRow>
                <TableHead>Active</TableHead>
                <TableHead>User Agent</TableHead>
                <TableHead>IP</TableHead>
                <TableHead>Last IP</TableHead>
                <TableHead>Last Used At</TableHead>
                <TableHead>Created At</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.sessions}>
                {(session) => (
                  <TableRow>
                    <TableCell>
                      <Show when={session.active} fallback={<div class="mx-auto h-4 w-4 rounded-full bg-gray-500" title="Inactive" />}>
                        <div class="mx-auto h-4 w-4 rounded-full bg-green-500" title="Active" />
                      </Show>
                    </TableCell>
                    <TableCell>{session.userAgent}</TableCell>
                    <TableCell>{session.ip}</TableCell>
                    <TableCell>{session.lastIp}</TableCell>
                    <TableCell>{formatDate(parseDate(session.lastUsedAtTime))}</TableCell>
                    <TableCell>{formatDate(parseDate(session.createdAtTime))}</TableCell>
                    <TableCell class="py-0">
                      <Show when={!session.current} fallback={<Badge>Current</Badge>}>
                        <Button variant="destructive" size="sm" onClick={() => setRevokeMySessionsConfirm(session.id)}>
                          Revoke
                        </Button>
                      </Show>
                    </TableCell>
                  </TableRow>
                )}
              </For>
            </TableBody>
          </TableRoot>
        </Suspense>

        <Shared.Title>Groups</Shared.Title>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <TableRoot>
            <TableCaption>{data()?.groups.length} {Humanize.pluralize(data()?.groups.length || 0, "Group")}</TableCaption>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Description</TableHead>
                <TableHead>Joined At</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.groups}>
                {(group) =>
                  <TableRow>
                    <TableCell>{group.name}</TableCell>
                    <TableCell>{group.description}</TableCell>
                    <TableCell>{formatDate(parseDate(group.joinedAtTime))}</TableCell>
                  </TableRow>
                }
              </For>
            </TableBody>
          </TableRoot>
        </Suspense>

      </ErrorBoundary>
    </LayoutNormal>
  )
}

type ChangeUsernameForm = {
  newUsername: string
}

const actionUpdateMyUsername = action((form: ChangeUsernameForm) => useClient()
  .user.updateMyUsername(form)
  .then(() => revalidate([getProfilePage.key, getSession.key]))
  .catch(throwAsFormError))

function ChangeUsernameForm() {
  const [changeUsernameForm, { Field, Form }] = createForm<ChangeUsernameForm>({ initialValues: { newUsername: "" } });
  const submit = useAction(actionUpdateMyUsername)

  return (
    <Form class="flex w-full max-w-xs flex-col gap-4" onSubmit={(form) => submit(form).then(() => reset(changeUsernameForm))}>
      <Field name="newUsername" validate={required("Please enter a new username.")}>
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>New username</FieldLabel>
            <FieldControl field={field}>
              <Input
                {...props}
                placeholder="New username"
                value={field.value}
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Button type="submit" disabled={changeUsernameForm.submitting}>
        <Show when={!changeUsernameForm.submitting} fallback={<>Updating username</>}>
          Update username
        </Show>
      </Button>
      <FormMessage form={changeUsernameForm} />
    </Form>
  )
}

type ChangePasswordForm = {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

const actionUpdateMyPassword = action((form: ChangePasswordForm) => useClient()
  .user.updateMyPassword(form)
  .then(() => revalidate(getProfilePage.key))
  .catch(throwAsFormError)
)

function ChangePasswordForm() {
  const [changePasswordForm, { Field, Form }] = createForm<ChangePasswordForm>({
    initialValues: { oldPassword: "", newPassword: "", confirmPassword: "" },
    validate: (form) => {
      if (form.newPassword != form.confirmPassword) {
        return {
          confirmPassword: "Password does not match."
        }
      }
      return {}
    }
  });
  const submit = useAction(actionUpdateMyPassword)

  return (
    <Form class="flex w-full max-w-xs flex-col gap-4" onSubmit={(form) => submit(form).then(() => reset(changePasswordForm))}>
      <input class="hidden" type="text" name="username" autocomplete="username" />
      <Field name="oldPassword" validate={required("Please enter your old password.")}>
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>Old password</FieldLabel>
            <FieldControl field={field}>
              <Input
                {...props}
                autocomplete="current-password"
                placeholder="Old password"
                type="password"
                value={field.value}
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="newPassword" validate={required("Please enter a new password.")}>
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>New password</FieldLabel>
            <FieldControl field={field}>
              <Input
                {...props}
                autocomplete="new-password"
                placeholder="New password"
                type="password"
                value={field.value}
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="confirmPassword">
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>Confirm new password</FieldLabel>
            <FieldControl field={field}>
              <Input
                {...props}
                autocomplete="new-password"
                placeholder="Confirm new password"
                type="password"
                value={field.value}
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Button type="submit" disabled={changePasswordForm.submitting}>
        <Show when={changePasswordForm.submitting} fallback={<>Update password</>}>
          Updating password
        </Show>
      </Button>
      <FormMessage form={changePasswordForm} />
    </Form>
  )
}

function Center(props: ParentProps) {
  return (
    <div class="flex justify-center">
      {props.children}
    </div>
  )
}

