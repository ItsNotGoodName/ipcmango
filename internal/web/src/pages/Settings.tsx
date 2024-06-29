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
import { Label } from "~/ui/Label";

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

export default function () {
  const data = createQuery(() => api.settings.get);

  const syncVideoInModeMutation = useMutation();
  const locationMutation = useMutation();

  const sunOffsetMutation = useMutation();
  const [sunOffsetEdit, setSunOffsetEdit] = createSignal(false);
  const [sunsetOffset, setSunsetOffset] = createSignal("");
  const [sunriseOffrise, setSunriseOffset] = createSignal("");
  const showSunOffsetForm = () => {
    if (sunOffsetEdit() == true) return;
    batch(() => {
      setSunOffsetEdit(true);
      setSunriseOffset(data.data?.sunrise_offset || "");
      setSunsetOffset(data.data?.sunset_offset || "");
    });
  };
  const submitSunOffset = () =>
    sunOffsetMutation.mutate(
      {
        sunrise_offset: sunriseOffrise(),
        sunset_offset: sunsetOffset(),
      },
      {
        onSuccess: () => setSunOffsetEdit(false),
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

      <div class="flex flex-col gap-2">
        <div class="flex gap-2">
          <TextFieldRoot
            disabled={sunOffsetMutation.isPending}
            validationState={sunOffsetMutation.error ? "invalid" : "valid"}
            onFocusIn={showSunOffsetForm}
            value={
              sunOffsetEdit() ? sunriseOffrise() : data.data?.sunrise_offset
            }
            onChange={setSunriseOffset}
            class="flex-1 space-y-2"
          >
            <TextFieldLabel>Sunrise Offset</TextFieldLabel>
            <TextFieldInput></TextFieldInput>
            <TextFieldDescription>TODO</TextFieldDescription>
            <TextFieldErrorMessage>
              {sunOffsetMutation.error?.message}
            </TextFieldErrorMessage>
          </TextFieldRoot>
          <TextFieldRoot
            disabled={sunOffsetMutation.isPending}
            validationState={sunOffsetMutation.error ? "invalid" : "valid"}
            onFocusIn={showSunOffsetForm}
            value={sunOffsetEdit() ? sunsetOffset() : data.data?.sunset_offset}
            onChange={setSunsetOffset}
            class="flex-1 space-y-2"
          >
            <TextFieldLabel>Sunset Offset</TextFieldLabel>
            <TextFieldInput></TextFieldInput>
            <TextFieldDescription>TODO</TextFieldDescription>
            <TextFieldErrorMessage>
              {sunOffsetMutation.error?.message}
            </TextFieldErrorMessage>
          </TextFieldRoot>
        </div>
        <Show when={sunOffsetEdit()}>
          <div class="flex justify-end gap-2">
            <Button
              disabled={sunOffsetMutation.isPending}
              onClick={() => setSunOffsetEdit(false)}
            >
              Cancel
            </Button>
            <Button
              disabled={sunOffsetMutation.isPending}
              onClick={submitSunOffset}
            >
              Save
            </Button>
          </div>
        </Show>
      </div>
    </LayoutNormal>
  );
}
