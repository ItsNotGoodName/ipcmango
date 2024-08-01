import {
  createForm,
  reset,
  setValue,
  setValues,
  submit,
} from "@modular-forms/solid";
import { LayoutCenter } from "~/ui/Layout";
import {
  SwitchControl,
  SwitchErrorMessage,
  SwitchLabel,
  SwitchRoot,
} from "~/ui/Switch";
import { api } from "./data";
import {
  createMutation,
  createQuery,
  useQueryClient,
} from "@tanstack/solid-query";
import { PatchSettings, deleteApiSettings, patchApiSettings } from "~/client";
import {
  SelectContent,
  SelectErrorMessage,
  SelectItem,
  SelectLabel,
  SelectListbox,
  SelectPortal,
  SelectRoot,
  SelectTrigger,
  SelectValue,
} from "~/ui/Select";
import {
  TextFieldErrorMessage,
  TextFieldInput,
  TextFieldLabel,
  TextFieldRoot,
} from "~/ui/TextField";
import { ErrorBoundary, Show, Suspense, batch, createSignal } from "solid-js";
import { Button } from "~/ui/Button";
import { Seperator } from "~/ui/Seperator";
import { validationState } from "~/lib/utils";
import {
  NumberFieldErrorMessage,
  NumberFieldInput,
  NumberFieldLabel,
  NumberFieldRoot,
} from "~/ui/NumberField";
import { PageError, PageTitle } from "~/ui/Page";
import { toast } from "~/ui/Toast";
import {
  AlertDialogRoot,
  AlertDialogTrigger,
  AlertDialogTitle,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogModal,
} from "~/ui/AlertDialog";
import { Skeleton } from "~/ui/Skeleton";
import { AlertDescription, AlertRoot } from "~/ui/Alert";
import { useBeforeLeave } from "@solidjs/router";

type Form = {
  latitude?: number | string;
  location?: string;
  longitude?: number | string;
  sunrise_offset?: string;
  sunset_offset?: string;
  sync_video_in_mode?: boolean;
};

export default function Settings() {
  const client = useQueryClient();

  const data = createQuery(() => ({
    ...api.settings.get,
    throwOnError: true,
  }));

  const locations = createQuery(() => ({
    ...api.locations.list,
    throwOnError: true,
  }));

  // Form
  const formMutation = createMutation(() => ({
    mutationFn: (requestBody: PatchSettings) =>
      patchApiSettings({ requestBody }),
    onSuccess: (data) => client.setQueryData(api.settings.get.queryKey, data),
  }));
  const [form, { Field, Form }] = createForm<Form>();
  const submitForm = (value: Form) =>
    formMutation.mutateAsync(
      {
        ...value,
        latitude:
          value.latitude != undefined ? Number(value.latitude) : undefined,
        longitude:
          value.longitude != undefined ? Number(value.longitude) : undefined,
      },
      {
        onSuccess: () => reset(form),
      },
    );

  useBeforeLeave((e) => {
    if (form.dirty && !e.defaultPrevented) {
      // preventDefault to block immediately and prompt user async
      e.preventDefault();
      setTimeout(() => {
        if (window.confirm("Discard unsaved changes - are you sure?")) {
          // user wants to proceed anyway so retry with force=true
          e.retry(true);
        }
      }, 100);
    }
  });

  // Default
  const [defaultDialog, setDefaultDialog] = createSignal(false);
  const defaultMutation = createMutation(() => ({
    mutationFn: () => deleteApiSettings(),
    onSuccess: (data) =>
      batch(() => {
        client.setQueryData(api.settings.get.queryKey, data);
        setDefaultDialog(false);
        reset(form);
      }),
    onError: (error) => toast.error(error.name, error.message),
  }));

  // Detect coordinate
  const [coordinateDetectLoading, setCoordinateDetectLoading] =
    createSignal(false);
  const detectCoordinate = () => {
    if (coordinateDetectLoading()) return;
    setCoordinateDetectLoading(true);

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setCoordinateDetectLoading(false);
        setValues(form, {
          latitude: pos.coords.latitude,
          longitude: pos.coords.longitude,
        });
      },
      (error) =>
        batch(() => {
          setCoordinateDetectLoading(false);
          toast.error("Detect Coordinate", error.message);
        }),
      { timeout: 10000 },
    );
  };

  return (
    <LayoutCenter class="max-w-xl">
      <PageTitle>Settings</PageTitle>

      <ErrorBoundary fallback={(error) => <PageError error={error} />}>
        <Form onSubmit={submitForm} class="flex flex-col gap-4">
          <Field name="sync_video_in_mode" type="boolean">
            {(field, props) => (
              <SwitchRoot
                validationState={validationState(field.error)}
                checked={
                  form.dirty && field.value != undefined
                    ? field.value
                    : data.data?.sync_video_in_mode
                }
                onChange={(value) =>
                  setValue(form, "sync_video_in_mode", value)
                }
                class="flex items-center justify-between gap-2"
              >
                <div>
                  <SwitchLabel>Sync Video In Mode</SwitchLabel>
                  <SwitchErrorMessage>{field.error}</SwitchErrorMessage>
                </div>
                <SwitchControl {...props} />
              </SwitchRoot>
            )}
          </Field>
        </Form>

        <Suspense fallback={<Skeleton class="h-32" />}>
          <Field name="location">
            {(field) => {
              const value = () =>
                form.dirty && field.value != undefined
                  ? field.value
                  : data.data?.location;
              return (
                <SelectRoot
                  validationState={validationState(field.error)}
                  options={locations.data || []}
                  value={value()}
                  onChange={(newValue) =>
                    newValue != value() && setValue(form, "location", newValue)
                  }
                  itemComponent={(props) => (
                    <SelectItem item={props.item}>
                      {props.item.rawValue}
                    </SelectItem>
                  )}
                  class="space-y-2"
                >
                  <SelectLabel>Location</SelectLabel>
                  <SelectTrigger>
                    <SelectValue<string>>
                      {(state) => state.selectedOption()}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectPortal>
                    <SelectContent>
                      <SelectListbox />
                    </SelectContent>
                  </SelectPortal>
                  <SelectErrorMessage>{field.error}</SelectErrorMessage>
                </SelectRoot>
              );
            }}
          </Field>
        </Suspense>

        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-4 sm:flex-row">
              <Field name="sunrise_offset">
                {(field, props) => (
                  <TextFieldRoot
                    validationState={validationState(field.error)}
                    value={
                      field.dirty && field.value != undefined
                        ? field.value
                        : data.data?.sunrise_offset
                    }
                    class="flex-1 space-y-2"
                  >
                    <TextFieldLabel>Sunrise Offset</TextFieldLabel>
                    <TextFieldInput {...props} />
                    <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
                  </TextFieldRoot>
                )}
              </Field>
              <Field name="sunset_offset">
                {(field, props) => (
                  <TextFieldRoot
                    validationState={validationState(field.error)}
                    value={
                      field.dirty && field.value != undefined
                        ? field.value
                        : data.data?.sunset_offset
                    }
                    class="flex-1 space-y-2"
                  >
                    <TextFieldLabel>Sunset Offset</TextFieldLabel>
                    <TextFieldInput {...props} />
                    <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
                  </TextFieldRoot>
                )}
              </Field>
            </div>
          </div>
        </Suspense>

        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-4 sm:flex-row">
              <Field name="latitude" type="number">
                {(field) => (
                  <NumberFieldRoot
                    validationState={validationState(field.error)}
                    value={
                      form.dirty && field.value != undefined
                        ? field.value
                        : data.data?.latitude
                    }
                    onChange={(value) => setValue(form, "latitude", value)}
                    class="flex-1 space-y-2"
                  >
                    <NumberFieldLabel>Latitude</NumberFieldLabel>
                    <NumberFieldInput />
                    <NumberFieldErrorMessage>
                      {field.error}
                    </NumberFieldErrorMessage>
                  </NumberFieldRoot>
                )}
              </Field>
              <Field name="longitude" type="number">
                {(field) => (
                  <NumberFieldRoot
                    validationState={validationState(field.error)}
                    value={
                      form.dirty && field.value != undefined
                        ? field.value
                        : data.data?.longitude
                    }
                    onChange={(value) => setValue(form, "longitude", value)}
                    class="flex-1 space-y-2"
                  >
                    <NumberFieldLabel>Longitude</NumberFieldLabel>
                    <NumberFieldInput />
                    <NumberFieldErrorMessage>
                      {field.error}
                    </NumberFieldErrorMessage>
                  </NumberFieldRoot>
                )}
              </Field>
            </div>
            <div class="flex flex-col gap-2 sm:flex-row-reverse">
              <Button
                disabled={coordinateDetectLoading()}
                onClick={detectCoordinate}
                type="button"
              >
                Detect
              </Button>
            </div>
          </div>
        </Suspense>
      </ErrorBoundary>

      <Seperator />

      <div class="flex justify-end gap-2 max-sm:flex-col">
        <Button
          type="submit"
          onClick={() => submit(form)}
          disabled={!form.dirty || form.submitting}
        >
          Submit
        </Button>
        <Button
          type="submit"
          onClick={() => reset(form)}
          disabled={!form.dirty}
          variant="secondary"
        >
          Cancel
        </Button>
        <AlertDialogRoot open={defaultDialog()} onOpenChange={setDefaultDialog}>
          <AlertDialogTrigger
            as={Button}
            disabled={defaultMutation.isPending}
            variant="destructive"
          >
            Default
          </AlertDialogTrigger>
          <AlertDialogModal>
            <AlertDialogHeader>
              <AlertDialogTitle>
                Are you sure you wish to default settings?
              </AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                disabled={defaultMutation.isPending}
                onClick={() => defaultMutation.mutate()}
                variant="destructive"
              >
                Default
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogModal>
        </AlertDialogRoot>
      </div>

      <Show when={form.response.message}>
        <AlertRoot variant="destructive">
          <AlertDescription>{form.response.message}</AlertDescription>
        </AlertRoot>
      </Show>
    </LayoutCenter>
  );
}
