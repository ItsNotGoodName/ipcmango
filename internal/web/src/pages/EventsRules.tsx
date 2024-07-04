import Humanize from "humanize-plus";
import { A } from "@solidjs/router";
import {
  createMutation,
  createQuery,
  useQueryClient,
} from "@tanstack/solid-query";
import { RiSystemAddLine, RiSystemDeleteBinLine } from "solid-icons/ri";
import { ErrorBoundary, For, Suspense, createSignal } from "solid-js";
import {
  createRowSelection,
  createValueDialog,
  validationState,
} from "~/lib/utils";
import {
  BreadcrumbsRoot,
  BreadcrumbsItem,
  BreadcrumbsLink,
  BreadcrumbsSeparator,
} from "~/ui/Breadcrumbs";
import { Button } from "~/ui/Button";
import {
  CheckboxRoot,
  CheckboxControl,
  CheckboxErrorMessage,
  CheckboxLabel,
} from "~/ui/Checkbox";
import { LayoutCenter } from "~/ui/Layout";
import { PageError, PageTitle } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import {
  TableBody,
  TableCell,
  TableHead,
  TableHeadBase,
  TableHeader,
  TableRoot,
  TableRow,
} from "~/ui/Table";
import { api } from "./data";
import {
  CreateEventRule,
  UpdateEventRule,
  deleteApiEventRules,
  postApiEventRulesByUuid,
  postApiEventRulesCreate,
} from "~/client";
import { toast } from "~/ui/Toast";
import {
  AlertDialogRoot,
  AlertDialogModal,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogCancel,
  AlertDialogAction,
  AlertDialogTrigger,
} from "~/ui/AlertDialog";
import {
  DialogRoot,
  DialogTrigger,
  DialogPortal,
  DialogOverlay,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogOverflow,
} from "~/ui/Dialog";
import { createForm, reset, setValue } from "@modular-forms/solid";
import {
  TextFieldErrorMessage,
  TextFieldInput,
  TextFieldLabel,
  TextFieldRoot,
} from "~/ui/TextField";
import { FormMessage } from "~/ui/Form";

export default function EventsRules() {
  const client = useQueryClient();

  const data = createQuery(() => ({
    ...api.eventRules.list,
    throwOnError: true,
  }));

  const rowSelection = createRowSelection(
    () =>
      data.data?.map((v) => ({ id: v.uuid, disabled: !v.can_delete })) ?? [],
  );

  const deleteModal = createValueDialog<Array<{ id: string; code: string }>>(
    [],
  );
  const deleteMutation = createMutation(() => ({
    mutationFn: () =>
      deleteApiEventRules({ requestBody: rowSelection.selections() }),
    onSuccess: () =>
      client
        .invalidateQueries({ queryKey: api.eventRules.list.queryKey })
        .then(deleteModal.setClose),
    onError: (error) => toast.error(error.name, error.message),
  }));

  const [createModal, setCreateModal] = createSignal(false);

  return (
    <LayoutCenter class="max-w-4xl">
      <PageTitle>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            <BreadcrumbsLink as={A} href="/events">
              Events
            </BreadcrumbsLink>
            <BreadcrumbsSeparator />
          </BreadcrumbsItem>
          <BreadcrumbsItem>Rules</BreadcrumbsItem>
        </BreadcrumbsRoot>
      </PageTitle>

      <div class="flex justify-end gap-2">
        <DialogRoot open={createModal()} onOpenChange={setCreateModal}>
          <DialogTrigger as={Button} size="icon">
            <RiSystemAddLine class="size-5" />
          </DialogTrigger>
          <DialogPortal>
            <DialogOverlay />
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Create Event Rule</DialogTitle>
              </DialogHeader>
              <DialogOverflow>
                <CreateForm onClose={() => setCreateModal(false)} />
              </DialogOverflow>
            </DialogContent>
          </DialogPortal>
        </DialogRoot>

        <AlertDialogRoot
          open={deleteModal.open()}
          onOpenChange={(value) =>
            value
              ? deleteModal.setValue(
                  rowSelection.items
                    .map((id, index) => ({
                      ...id,
                      code: data.data![index].code || "",
                    }))
                    .filter((v) => v.checked),
                )
              : deleteModal.setClose()
          }
        >
          <AlertDialogTrigger
            as={Button}
            size="icon"
            variant="destructive"
            disabled={!rowSelection.multiple() || deleteMutation.isPending}
          >
            <RiSystemDeleteBinLine class="size-5" />
          </AlertDialogTrigger>
          <AlertDialogModal>
            <AlertDialogHeader>
              <AlertDialogTitle>
                Are you sure you wish to delete {deleteModal.value().length}{" "}
                event{" "}
                {Humanize.pluralize(
                  deleteModal.value().length,
                  "rule",
                  "rules",
                )}
                ?
              </AlertDialogTitle>
              <AlertDialogDescription>
                <ul>
                  <For each={deleteModal.value()}>
                    {(e) => <li>{e.code}</li>}
                  </For>
                </ul>
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                variant="destructive"
                disabled={deleteMutation.isPending}
                onClick={() => deleteMutation.mutate()}
              >
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogModal>
        </AlertDialogRoot>
      </div>

      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <TableRoot>
            <TableHeader>
              <TableRow>
                <TableHead>
                  <CheckboxRoot
                    indeterminate={rowSelection.multiple()}
                    checked={rowSelection.all()}
                    onChange={rowSelection.setAll}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableHead>
                <TableHead class="w-full">Code</TableHead>
                <TableHeadBase>DB</TableHeadBase>
                <TableHeadBase>Live</TableHeadBase>
                <TableHeadBase>MQTT</TableHeadBase>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data.data}>
                {(item, index) => {
                  const update = createMutation(() => ({
                    mutationFn: (requestBody: UpdateEventRule) =>
                      postApiEventRulesByUuid({ uuid: item.uuid, requestBody }),
                    onSuccess: () =>
                      void client.invalidateQueries({
                        queryKey: api.eventRules.list.queryKey,
                      }),
                  }));
                  const patch = (data: any) =>
                    update.mutate({
                      code: item.code,
                      allow_db: item.allow_db,
                      allow_live: item.allow_live,
                      allow_mqtt: item.allow_mqtt,
                      ...data,
                    });

                  return (
                    <TableRow>
                      <TableCell>
                        <CheckboxRoot
                          disabled={rowSelection.items[index()]?.disabled}
                          checked={rowSelection.items[index()]?.checked}
                          onChange={(value) =>
                            rowSelection.set(item.uuid, value)
                          }
                        >
                          <CheckboxControl />
                        </CheckboxRoot>
                      </TableCell>
                      <TableCell class="w-full">{item.code}</TableCell>
                      <TableCell>
                        <div class="flex justify-center">
                          <CheckboxRoot
                            checked={item.allow_db}
                            onChange={(allow_db) => patch({ allow_db })}
                          >
                            <CheckboxControl />
                          </CheckboxRoot>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div class="flex justify-center">
                          <CheckboxRoot
                            checked={item.allow_live}
                            onChange={(allow_live) => patch({ allow_live })}
                          >
                            <CheckboxControl />
                          </CheckboxRoot>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div class="flex justify-center">
                          <CheckboxRoot
                            checked={item.allow_mqtt}
                            onChange={(allow_mqtt) => patch({ allow_mqtt })}
                          >
                            <CheckboxControl />
                          </CheckboxRoot>
                        </div>
                      </TableCell>
                    </TableRow>
                  );
                }}
              </For>
            </TableBody>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutCenter>
  );
}

function CreateForm(props: { onClose: () => void }) {
  const [createMore, setCreateMore] = createSignal(false);
  const client = useQueryClient();
  const [form, { Field, Form }] = createForm<CreateEventRule>({
    initialValues: {
      allow_db: false,
      allow_live: false,
      allow_mqtt: false,
      code: "",
    },
  });
  const mutation = createMutation(() => ({
    mutationFn: (requestBody: CreateEventRule) =>
      postApiEventRulesCreate({ requestBody }),
    onSuccess: (_, req) =>
      client
        .invalidateQueries({ queryKey: api.eventRules.list.queryKey })
        .then(() =>
          createMore()
            ? reset(form, {
                initialValues: {
                  code: "",
                  allow_db: req.allow_db,
                  allow_live: req.allow_live,
                  allow_mqtt: req.allow_mqtt,
                },
              })
            : props.onClose(),
        ),
  }));

  return (
    <Form onSubmit={(data) => mutation.mutateAsync(data)} class="space-y-4">
      <Field name="code">
        {(field, props) => (
          <TextFieldRoot
            validationState={validationState(field.error)}
            value={field.value}
            class="space-y-2"
          >
            <TextFieldLabel>Code</TextFieldLabel>
            <TextFieldInput {...props} />
            <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
          </TextFieldRoot>
        )}
      </Field>
      <div class="flex flex-wrap gap-4">
        <Field name="allow_db" type="boolean">
          {(field, props) => (
            <CheckboxRoot
              validationState={validationState(field.error)}
              checked={field.value}
              onChange={(value) => setValue(form, "allow_db", value)}
              class="flex justify-between gap-2"
            >
              <CheckboxLabel>Allow DB</CheckboxLabel>
              <CheckboxControl {...props} />
              <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
            </CheckboxRoot>
          )}
        </Field>
        <Field name="allow_live" type="boolean">
          {(field, props) => (
            <CheckboxRoot
              validationState={validationState(field.error)}
              checked={field.value}
              onChange={(value) => setValue(form, "allow_live", value)}
              class="flex justify-between gap-2"
            >
              <CheckboxLabel>Allow Live</CheckboxLabel>
              <CheckboxControl {...props} />
              <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
            </CheckboxRoot>
          )}
        </Field>
        <Field name="allow_mqtt" type="boolean">
          {(field, props) => (
            <CheckboxRoot
              validationState={validationState(field.error)}
              checked={field.value}
              onChange={(value) => setValue(form, "allow_mqtt", value)}
              class="flex justify-between gap-2"
            >
              <CheckboxLabel>Allow MQTT</CheckboxLabel>
              <CheckboxControl {...props} />
              <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
            </CheckboxRoot>
          )}
        </Field>
      </div>
      <Button disabled={form.submitting} class="w-full">
        Create
      </Button>
      <FormMessage form={form} />
      <CheckboxRoot
        checked={createMore()}
        onChange={setCreateMore}
        class="flex gap-2"
      >
        <CheckboxLabel>Add More</CheckboxLabel>
        <CheckboxControl {...props} />
      </CheckboxRoot>
    </Form>
  );
}
