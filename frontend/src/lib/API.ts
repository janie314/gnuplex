interface MediaDir {
  Path: string;
  LastScanned: string;
}

interface MediaItem {
  ID: number;
  Path: string;
  LastPlayed: string;
}

interface MediaItemRes {
  res: MediaItem[];
  count: number;
}

interface FileExtension {
  ID: number;
  Extension: string;
}

interface SubTrack {
  id: number;
  title: string;
  selected: boolean;
}

interface Version {
  version: string;
  source_hash: string;
}

const headers = { "Content-Type": "application/json" };
const post = { method: "POST" };
const json = (res: Response) => res.json();

class API {
  public static async play() {
    return await fetch("/api/play", post);
  }

  public static async pause() {
    return await fetch("/api/pause", post);
  }

  public static async playpause() {
    return await fetch("/api/playpause", post);
  }

  public static async skip() {
    return await fetch("/api/skip", post);
  }

  public static async getPos(): Promise<number> {
    return await fetch("/api/pos").then(json);
  }

  public static async getTimeRemaining(): Promise<number> {
    return await fetch("/api/timeremaining").then(json);
  }

  public static async getVersion(): Promise<Version> {
    return await fetch("/api/version").then(json);
  }

  public static async setPos(pos: number) {
    return await fetch("/api/pos", {
      ...post,
      headers,
      body: JSON.stringify({ pos }),
    }).then(json);
  }

  public static async getPaused(): Promise<boolean> {
    return await fetch("/api/paused").then(json);
  }

  public static async getVol(): Promise<number> {
    return await fetch("/api/vol").then(json);
  }

  public static async setVol(vol: number) {
    return await fetch("/api/vol", {
      ...post,
      headers,
      body: JSON.stringify({ vol }),
    }).then(json);
  }

  public static async getNowPlaying(): Promise<MediaItem[] | null> {
    return await fetch("/api/nowplaying").then(json);
  }

  public static async playMedia(
    mediaItem: MediaItem,
    play_next: boolean,
    play_last: boolean,
  ) {
    return await fetch("/api/playmedia", {
      ...post,
      headers,
      body: JSON.stringify({ id: mediaItem.ID, play_next, play_last }),
    });
  }

  public static async cast(url: string, temp: boolean) {
    return await fetch("/api/cast", {
      ...post,
      headers,
      body: JSON.stringify({ url, temp }),
    });
  }

  public static async getMediadirs(): Promise<MediaDir[]> {
    return await fetch("/api/mediadirs").then(json);
  }

  public static async setMediadirs(mediadirs: string[]) {
    return await fetch("/api/mediadirs", {
      ...post,
      headers,
      body: JSON.stringify(mediadirs),
    });
  }

  public static async getFileExts() {
    return await fetch("/api/file_exts")
      .then(json)
      .then((data: FileExtension[]) => {
        return data.sort((a, b) =>
          a.Extension.toLowerCase() < b.Extension.toLowerCase() ? -1 : 1,
        );
      });
  }

  public static async setFileExts(file_exts: string[]) {
    return await fetch("/api/file_exts", {
      ...post,
      headers,
      body: JSON.stringify(file_exts),
    });
  }

  public static async getMediaItems(
    search: string,
    paginationOffset: number,
  ): Promise<MediaItemRes> {
    const param = search || "";
    return await fetch(
      `/api/mediaitems?search=${encodeURIComponent(param)}&offset=${encodeURIComponent(paginationOffset)}`,
    ).then(json);
  }

  public static async scanLib() {
    return await fetch("/api/scanlib", post);
  }

  public static async getLast25Played(): Promise<MediaItem[]> {
    return await fetch("/api/last25").then(json);
  }

  public static async getSubTracks(): Promise<SubTrack[] | null> {
    return await fetch("/api/sub").then(json);
  }

  public static async setSubTrack(id: number) {
    const body = id === -1 ? { visible: false } : { visible: true, id };
    return await fetch("/api/sub", {
      ...post,
      headers,
      body: JSON.stringify(body),
    });
  }

  public static async setFilter(filter: string) {
    return await fetch("/api/filter", {
      ...post,
      headers,
      body: JSON.stringify({ filter }),
    });
  }
}

export { type MediaDir, type MediaItem, type SubTrack, API };
