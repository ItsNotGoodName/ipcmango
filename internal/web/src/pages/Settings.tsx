import { createForm, reset, setValue } from "@modular-forms/solid";
import { LayoutNormal } from "~/ui/Layout";
import {
  SwitchControl,
  SwitchDescription,
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
import { PatchSettings, patchApiSettings } from "~/client";
import {
  SelectContent,
  SelectDescription,
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
  TextFieldDescription,
  TextFieldErrorMessage,
  TextFieldInput,
  TextFieldLabel,
  TextFieldRoot,
} from "~/ui/TextField";
import { Show, batch, createSignal } from "solid-js";
import { Button } from "~/ui/Button";
import { Seperator } from "~/ui/Seperator";
import { FormMessage } from "~/ui/Form";
import { validationState } from "~/lib/utils";
import {
  NumberFieldDescription,
  NumberFieldErrorMessage,
  NumberFieldInput,
  NumberFieldLabel,
  NumberFieldRoot,
} from "~/ui/NumberField";

function useMutation() {
  const client = useQueryClient();
  return createMutation(() => ({
    mutationFn: (requestBody: PatchSettings) =>
      patchApiSettings({
        requestBody,
      }),
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

function createToggle(formOpen: boolean, onOpen: () => void) {
  const [open, setOpen] = createSignal(formOpen);
  return {
    open,
    setOpen: () => {
      if (open() == true) return;
      batch(() => {
        setOpen(true);
        onOpen();
      });
    },
    setClose: () => setOpen(false),
  };
}

export default function () {
  const data = createQuery(() => api.settings.get);

  const syncVideoInModeMutation = useMutation();
  const locationMutation = useMutation();

  const [sunForm, { Field: SunField, Form: SunForm }] =
    createForm<UpdateSunForm>({
      initialValues: {
        sunrise_offset: "",
        sunset_offset: "",
      },
    });
  const sunFormToggle = createToggle(false, () => {
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

  const [coordinateForm, { Field: CoordinateField, Form: CoordinateForm }] =
    createForm<UpdateCoordinateForm>({
      initialValues: {
        latitude: 0,
        longitude: 0,
      },
    });
  const coordinateFormToggle = createToggle(false, () => {
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

  return (
    <LayoutNormal class="max-w-lg">
      <h1>Settings</h1>

      <SwitchRoot
        validationState={syncVideoInModeMutation.error ? "invalid" : "valid"}
        disabled={syncVideoInModeMutation.isPending}
        checked={data.data?.sync_video_in_mode}
        class="flex items-center justify-between gap-2"
        onChange={(isChecked) =>
          syncVideoInModeMutation.mutate({ sync_video_in_mode: isChecked })
        }
      >
        <div>
          <SwitchLabel>Sync Video In Mode</SwitchLabel>
          <SwitchDescription>TODO</SwitchDescription>
          <SwitchErrorMessage>
            {syncVideoInModeMutation.error?.message}
          </SwitchErrorMessage>
        </div>
        <SwitchControl />
      </SwitchRoot>

      <Seperator />

      <SelectRoot
        validationState={locationMutation.error ? "invalid" : "valid"}
        disabled={locationMutation.isPending}
        options={["Local", "UTC"]}
        value={data.data?.location}
        onChange={(value) => locationMutation.mutate({ location: value })}
        itemComponent={(props) => (
          <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
        )}
        class="space-y-2"
      >
        <SelectLabel>Location</SelectLabel>
        <SelectTrigger>
          <SelectValue<string>>{(state) => state.selectedOption()}</SelectValue>
        </SelectTrigger>
        <SelectPortal>
          <SelectContent>
            <SelectListbox />
          </SelectContent>
        </SelectPortal>
        <SelectDescription>TODO</SelectDescription>
        <SelectErrorMessage>
          {locationMutation.error?.message}
        </SelectErrorMessage>
      </SelectRoot>

      <Seperator />

      <SunForm onSubmit={submitSunForm} class="flex flex-col gap-2">
        <div class="flex gap-4">
          <SunField name="sunrise_offset">
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                onFocusIn={sunFormToggle.setOpen}
                value={
                  sunFormToggle.open() ? field.value : data.data?.sunrise_offset
                }
                class="flex-1 space-y-2"
              >
                <TextFieldLabel>Sunrise Offset</TextFieldLabel>
                <TextFieldInput {...props} />
                <TextFieldDescription>TODO</TextFieldDescription>
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
                  sunFormToggle.open() ? field.value : data.data?.sunset_offset
                }
                class="flex-1 space-y-2"
              >
                <TextFieldLabel>Sunset Offset</TextFieldLabel>
                <TextFieldInput {...props} />
                <TextFieldDescription>TODO</TextFieldDescription>
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </SunField>
        </div>
        <Show when={sunFormToggle.open()}>
          <FormMessage form={sunForm} />
          <div class="flex justify-end gap-2">
            <Button
              disabled={sunForm.submitting}
              onClick={sunFormToggle.setClose}
              variant="secondary"
              type="button"
            >
              Cancel
            </Button>
            <Button disabled={sunForm.submitting} type="submit">
              Save
            </Button>
          </div>
        </Show>
      </SunForm>

      <Seperator />

      <CoordinateForm
        onSubmit={submitCoordinateForm}
        class="flex flex-col gap-2"
      >
        <div class="flex gap-4">
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
                <NumberFieldDescription>TODO</NumberFieldDescription>
                <NumberFieldErrorMessage>{field.error}</NumberFieldErrorMessage>
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
                <NumberFieldDescription>TODO</NumberFieldDescription>
                <NumberFieldErrorMessage>{field.error}</NumberFieldErrorMessage>
              </NumberFieldRoot>
            )}
          </CoordinateField>
        </div>
        <Show when={coordinateFormToggle.open()}>
          <FormMessage form={coordinateForm} />
          <div class="flex justify-end gap-2">
            <Button
              disabled={coordinateForm.submitting}
              onClick={coordinateFormToggle.setClose}
              variant="secondary"
              type="button"
            >
              Cancel
            </Button>
            <Button disabled={coordinateForm.submitting} type="submit">
              Save
            </Button>
          </div>
        </Show>
      </CoordinateForm>
    </LayoutNormal>
  );
}
