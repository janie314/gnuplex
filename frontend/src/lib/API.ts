interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

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

class API {
  public static async play() {
    return await fetch("/api/play", { method: "POST" });
  }

  public static async pause() {
    return await fetch("/api/pause", { method: "POST" });
  }

  public static async getPos(): Promise<number> {
    return await fetch("/api/pos")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get pos", e);
        return 0;
      });
  }

  public static async getTimeRemaining(): Promise<number> {
    return await fetch("/api/timeremaining")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get time remaining", e);
        return 0;
      });
  }

  public static async getVersion() {
    return (await fetch("/api/version").then((res) => res.json())) as string;
  }

  public static async setPos(pos: number) {
    return await fetch("/api/pos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ pos }),
    }).then((res) => res.json());
  }

  public static async getVol(): Promise<number> {
    return await fetch("/api/vol")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get vol", e);
        return 0;
      });
  }

  public static async setVol(vol: number) {
    return await fetch("/api/vol", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ vol }),
    }).then((res) => res.json());
  }

  public static async getNowPlaying() {
    return (await fetch("/api/nowplaying").then((res) =>
      res.json(),
    )) as MediaItem;
  }

  public static async setNowPlaying(mediaItem: MediaItem) {
    return await fetch("/api/nowplaying", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id: mediaItem.ID }),
    });
  }

  public static async cast(url: string, temp: boolean) {
    return await fetch("/api/cast", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ url, temp }),
    });
  }

  public static async getMediadirs(): Promise<MediaDir[]> {
    return await fetch("/api/mediadirs")
      .then((res) => res.json())
      .catch((e) => {
        console.error(e);
        return [];
      });
  }

  public static async setMediadirs(mediadirs: string[]) {
    return await fetch("/api/mediadirs", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(mediadirs),
    });
  }

  public static async getFileExts() {
    return await fetch("/api/file_exts")
      .then((res) => res.json())
      .then((data: FileExtension[]) => {
        return data.sort((a, b) =>
          a.Extension.toLowerCase() < b.Extension.toLowerCase() ? -1 : 1,
        );
      });
  }

  public static async setFileExts(file_exts: string[]) {
    return await fetch("/api/file_exts", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(file_exts),
    });
  }

  public static async getMediaItems(search: string, paginationOffset: number) {
    const param = search || "";
    return (await fetch(
      `/api/mediaitems?search=${encodeURIComponent(param)}&offset=${encodeURIComponent(paginationOffset)}`,
    ).then((res) => res.json())) as MediaItemRes;
  }

  public static async scanLib() {
    return await fetch("/api/scanlib", { method: "POST" });
  }

  public static async getLast25Played() {
    return (await fetch("/api/last25").then((res) =>
      res.json(),
    )) as MediaItem[];
  }

  public static sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  public static async cycleSub(next: boolean) {
    const dir = next ? "next" : "prev";
    return await fetch(`/api/sub?dir=${dir}`, { method: "POST" });
  }
}

export { type MediaDir, type MediaItem, API };
