import { createQuery } from "@tanstack/solid-query";
import { RiSystemFilterLine } from "solid-icons/ri";
import { Device } from "~/client";
import { api } from "~/pages/data";
import {
  ComboboxContent,
  ComboboxControl,
  ComboboxIcon,
  ComboboxInput,
  ComboboxItem,
  ComboboxItemLabel,
  ComboboxListbox,
  ComboboxReset,
  ComboboxRoot,
  ComboboxTrigger,
  ComboboxState,
} from "~/ui/Combobox";

export function DeviceFilterCombobox(props: {
  setDeviceIDs: (ids: string[]) => void;
  deviceIDs: string[];
}) {
  const data = createQuery(() => api.devices.list);

  return (
    <ComboboxRoot<Device>
      multiple
      optionValue="uuid"
      optionTextValue="name"
      optionLabel="name"
      options={data.data || []}
      placeholder="Device"
      value={data.data?.filter((v) => props.deviceIDs.includes(v.uuid))}
      onChange={(value) => props.setDeviceIDs(value.map((v) => v.uuid))}
      itemComponent={(props) => (
        <ComboboxItem item={props.item}>
          <ComboboxItemLabel>{props.item.rawValue.name}</ComboboxItemLabel>
        </ComboboxItem>
      )}
    >
      <ComboboxControl<Device> aria-label="Device">
        {(state) => (
          <ComboboxTrigger>
            <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
            Device
            <ComboboxState
              state={state}
              getOptionString={(option) => option.name}
            />
            <ComboboxReset state={state} class="size-4" />
          </ComboboxTrigger>
        )}
      </ComboboxControl>
      <ComboboxContent>
        <ComboboxInput />
        <ComboboxListbox />
      </ComboboxContent>
    </ComboboxRoot>
  );
}
