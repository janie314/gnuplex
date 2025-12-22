import type { LongPressFns } from "@uidotdev/usehooks";
import type { MediaItem } from "../lib/API";
import { PageSelector } from "./Medialist/PageSelector";

function Medialist(props: {
  mediaItems: (MediaItem | null)[];
  subtitle: string;
  mediaItemCount: number | null;
  paginationOffset: number | null;
  setPaginationOffset: React.Dispatch<React.SetStateAction<number>> | null;
  longPressHook: LongPressFns;
}) {
  if (props.mediaItems.length === 0 || props.mediaItems[0] === null) {
    return null;
  }

  return (
    <div className="w-full flex flex-col mb-2 pl-2 whitespace-pre=wrap">
      <div className="flex flex-row align-center mb-1">
        <h1 className="header">{props.subtitle}</h1>
        <PageSelector
          mediaItemCount={props.mediaItemCount}
          paginationOffset={props.paginationOffset}
          setPaginationOffset={props.setPaginationOffset}
        />
      </div>
      {props.mediaItems
        .filter((mediaItem) => mediaItem !== null)
        .map((mediaItem, _i: number) => (
          <input
            type="button"
            className="inline-block w-screen max-w-screen-md text-left hover:bg-cyan-200 hover:dark:bg-cyan-700 p-1 border-lightgray dark:border-stone-800 border-t-2 whitespace-break-spaces break-words dark:bg-stone-950 dark:text-slate-200"
            key={props.subtitle + mediaItem.ID}
            value={mediaItem.Path}
            {...props.longPressHook}
          />
        ))}
    </div>
  );
}

export { Medialist };
