import { createForm, reset, setValue, setValues } from "@modular-forms/solid";
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
import { FormMessage } from "~/ui/Form";
import { createFormToggle, validationState } from "~/lib/utils";
import {
  NumberFieldErrorMessage,
  NumberFieldInput,
  NumberFieldLabel,
  NumberFieldRoot,
} from "~/ui/NumberField";
import { PageError, PageSubTitle, PageTitle } from "~/ui/Page";
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

function useMutation() {
  const client = useQueryClient();
  return createMutation(() => ({
    mutationFn: (requestBody: PatchSettings) =>
      patchApiSettings({ requestBody }),
    onSuccess: (data) => client.setQueryData(api.settings.get.queryKey, data),
  }));
}

type UpdateSunForm = {
  sunset_offset: string;
  sunrise_offset: string;
};

type UpdateCoordinateForm = {
  latitude: string | number;
  longitude: string | number;
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

  // Default
  const [defaultDialog, setDefaultDialog] = createSignal(false);
  const defaultMutation = createMutation(() => ({
    mutationFn: () => deleteApiSettings(),
    onSuccess: (data) =>
      batch(() => {
        client.setQueryData(api.settings.get.queryKey, data);
        setDefaultDialog(false);
      }),
    onError: (error) => toast.error(error.name, error.message),
  }));

  const syncVideoInModeMutation = useMutation();
  const locationMutation = useMutation();

  // Sun
  const [sunForm, { Field: SunField, Form: SunForm }] =
    createForm<UpdateSunForm>({
      initialValues: {
        sunrise_offset: "",
        sunset_offset: "",
      },
    });
  const sunFormToggle = createFormToggle(false, () => {
    reset(sunForm, {
      initialValues: {
        sunrise_offset: data.data?.sunrise_offset || "",
        sunset_offset: data.data?.sunset_offset || "",
      },
    });
  });
  const sunFormMutation = useMutation();
  const submitSunForm = (value: UpdateSunForm) =>
    sunFormMutation.mutateAsync(value, {
      onSuccess: sunFormToggle.setClose,
    });

  // Coordinate
  const [coordinateForm, { Field: CoordinateField, Form: CoordinateForm }] =
    createForm<UpdateCoordinateForm>({
      initialValues: {
        latitude: 0,
        longitude: 0,
      },
    });
  const coordinateFormToggle = createFormToggle(false, () => {
    reset(coordinateForm, {
      initialValues: {
        latitude: data.data?.latitude || 0,
        longitude: data.data?.longitude || 0,
      },
    });
  });
  const coordinateFormMutation = useMutation();
  const submitCoordinateForm = (value: UpdateCoordinateForm) =>
    coordinateFormMutation.mutateAsync(
      {
        latitude: Number(value.latitude),
        longitude: Number(value.longitude),
      },
      {
        onSuccess: coordinateFormToggle.setClose,
      },
    );
  const [coordinateDetectLoading, setCoordinateDetectLoading] =
    createSignal(false);
  const detectCoordinate = () => {
    if (coordinateDetectLoading()) return;
    setCoordinateDetectLoading(true);

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        batch(() => {
          setCoordinateDetectLoading(false);
          coordinateFormToggle.setOpen();
        });
        setValues(coordinateForm, {
          latitude: pos.coords.latitude,
          longitude: pos.coords.longitude,
        });
      },
      (error) =>
        batch(() => {
          setCoordinateDetectLoading(false);
          toast.error("Error", error.message);
        }),
      { timeout: 10000 },
    );
  };

  return (
    <LayoutCenter class="max-w-xl">
      <PageTitle>Settings</PageTitle>

      <ErrorBoundary fallback={(error) => <PageError error={error} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <SwitchRoot
            validationState={validationState(syncVideoInModeMutation.error)}
            disabled={syncVideoInModeMutation.isPending}
            checked={data.data?.sync_video_in_mode}
            class="flex items-center justify-between gap-2"
            onChange={(isChecked) =>
              syncVideoInModeMutation.mutate({ sync_video_in_mode: isChecked })
            }
          >
            <div>
              <SwitchLabel>Sync Video In Mode</SwitchLabel>
              <SwitchErrorMessage>
                {syncVideoInModeMutation.error?.message}
              </SwitchErrorMessage>
            </div>
            <SwitchControl />
          </SwitchRoot>
        </Suspense>

        <Seperator />

        <Suspense fallback={<Skeleton class="h-32" />}>
          <SelectRoot
            validationState={locationMutation.error ? "invalid" : "valid"}
            disabled={locationMutation.isPending}
            options={locations.data || []}
            value={data.data?.location}
            onChange={(value) => locationMutation.mutate({ location: value })}
            itemComponent={(props) => (
              <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
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
            <SelectErrorMessage>
              {locationMutation.error?.message}
            </SelectErrorMessage>
          </SelectRoot>
        </Suspense>

        <PageSubTitle>Sun Offset</PageSubTitle>

        <Suspense fallback={<Skeleton class="h-32" />}>
          <SunForm onSubmit={submitSunForm} class="flex flex-col gap-4">
            <div class="flex flex-col gap-4 sm:flex-row">
              <SunField name="sunrise_offset">
                {(field, props) => (
                  <TextFieldRoot
                    validationState={validationState(field.error)}
                    onFocusIn={sunFormToggle.setOpen}
                    value={
                      sunFormToggle.open()
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
              </SunField>
              <SunField name="sunset_offset">
                {(field, props) => (
                  <TextFieldRoot
                    validationState={validationState(field.error)}
                    onFocusIn={sunFormToggle.setOpen}
                    value={
                      sunFormToggle.open()
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
              </SunField>
            </div>
            <Show when={sunFormToggle.open()}>
              <div class="flex flex-col gap-2 sm:flex-row-reverse">
                <Button disabled={sunForm.submitting}>Save</Button>
                <Button
                  disabled={sunForm.submitting}
                  onClick={sunFormToggle.setClose}
                  variant="secondary"
                  type="button"
                >
                  Cancel
                </Button>
              </div>
            </Show>
            <FormMessage form={sunForm} />
          </SunForm>
        </Suspense>

        <PageSubTitle>Coordinate</PageSubTitle>

        <Suspense fallback={<Skeleton class="h-32" />}>
          <CoordinateForm
            onSubmit={submitCoordinateForm}
            class="flex flex-col gap-4"
          >
            <div class="flex flex-col gap-4 sm:flex-row">
              <CoordinateField name="latitude" type="number">
                {(field) => (
                  <NumberFieldRoot
                    validationState={validationState(field.error)}
                    onFocusIn={coordinateFormToggle.setOpen}
                    value={
                      coordinateFormToggle.open()
                        ? field.value
                        : data.data?.latitude
                    }
                    onChange={(value) =>
                      setValue(coordinateForm, "latitude", value)
                    }
                    class="flex-1 space-y-2"
                  >
                    <NumberFieldLabel>Latitude</NumberFieldLabel>
                    <NumberFieldInput />
                    <NumberFieldErrorMessage>
                      {field.error}
                    </NumberFieldErrorMessage>
                  </NumberFieldRoot>
                )}
              </CoordinateField>
              <CoordinateField name="longitude" type="number">
                {(field) => (
                  <NumberFieldRoot
                    validationState={validationState(field.error)}
                    onFocusIn={coordinateFormToggle.setOpen}
                    value={
                      coordinateFormToggle.open()
                        ? field.value
                        : data.data?.longitude
                    }
                    onChange={(value) =>
                      setValue(coordinateForm, "longitude", value)
                    }
                    class="flex-1 space-y-2"
                  >
                    <NumberFieldLabel>Longitude</NumberFieldLabel>
                    <NumberFieldInput />
                    <NumberFieldErrorMessage>
                      {field.error}
                    </NumberFieldErrorMessage>
                  </NumberFieldRoot>
                )}
              </CoordinateField>
            </div>
            <div class="flex flex-col gap-2 sm:flex-row-reverse">
              <Button
                disabled={coordinateDetectLoading()}
                onClick={detectCoordinate}
                type="button"
              >
                Detect
              </Button>
              <Show when={coordinateFormToggle.open()}>
                <Button disabled={coordinateForm.submitting}>Save</Button>
                <Button
                  disabled={coordinateForm.submitting}
                  onClick={coordinateFormToggle.setClose}
                  variant="secondary"
                  type="button"
                >
                  Cancel
                </Button>
              </Show>
            </div>
            <FormMessage form={coordinateForm} />
          </CoordinateForm>
        </Suspense>
      </ErrorBoundary>

      <Seperator />

      <div class="flex gap-2">
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
    </LayoutCenter>
  );
}
