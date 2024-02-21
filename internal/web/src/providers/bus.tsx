import { EventBus, EventHub, createEventBus, createEventHub } from '@solid-primitives/event-bus';
import {
  createContext,
  ParentComponent,
  useContext
} from "solid-js";
import { WSDahuaEvent } from '~/lib/models';

export type EventType = {
  action: string,
  data: unknown
}

type BusContextType = EventHub<{
  event: EventBus<EventType>;
  dahuaEvent: EventBus<WSDahuaEvent>
}>

const BusContext = createContext<BusContextType>();

type BusContextProps = {};

export const BusProvider: ParentComponent<BusContextProps> = (props) => {
  const store: BusContextType = createEventHub({
    event: createEventBus(),
    dahuaEvent: createEventBus()
  })

  return (
    <BusContext.Provider value={store}>
      {props.children}
    </BusContext.Provider>)
};

export function useBus(): BusContextType {
  const result = useContext(BusContext);
  if (!result) throw new Error("useBus must be used within a BusProvider");
  return result;
}
