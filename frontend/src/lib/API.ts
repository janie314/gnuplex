interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

class API {
  public static async play() {
    return await fetch("/api/play", { method: "POST" });
  }

  public static async pause() {
    return await fetch("/api/pause", { method: "POST" });
  }

  public static async getOriginPos(): Promise<number> {
    return await fetch("/api/pos")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get pos", e);
        return 0;
      });
  }

  public static async getOriginTimeRemaining() {
    return await fetch("/api/timeremaining")
      .then((res) => res.json())
      .then((res: IMPVRes) => {
        if (res.data !== undefined) {
          // @ts-ignore
          return Math.floor(res.data);
        }
        return 0;
      });
  }

  public static async getOriginVersion() {
    return (await fetch("/api/version").then((res) => res.json())) as string;
  }

  public static async setOriginPos(pos: number) {
    return await fetch(`/api/pos?pos=${pos}`, { method: "POST" }).then((res) =>
      res.json(),
    );
  }

  public static async getVol(): Promise<number> {
    return await fetch("/api/vol")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get vol", e);
        return 0;
      });
  }

  public static async setOriginVol(vol: number) {
    return await fetch(`/api/vol?vol=${vol}`, { method: "POST" }).then((res) =>
      res.json(),
    );
  }

  public static async getOriginMedia() {
    return await fetch("/api/media")
      .then((res) => res.json())
      .then((res: IMPVRes) => {
        if (res.data !== undefined) {
          return res.data as string;
        }
        return "";
      });
  }

  public static async setOriginMedia(mediafile: string) {
    return await fetch(`/api/media?mediafile=${encodeURI(mediafile)}`, {
      method: "POST",
    });
  }

  public static async getOriginMediadirs() {
    return await fetch("/api/mediadirs")
      .then((res) => res.json())
      .then((data: string[]) => {
        return data.sort((a, b) =>
          a.toLowerCase() < b.toLowerCase() ? -1 : 1,
        );
      });
  }

  public static async setOriginMediadirs(mediadirs: string[]) {
    return await fetch(
      `/api/mediadirs?mediadirs=${encodeURI(JSON.stringify(mediadirs))}`,
      { method: "POST" },
    );
  }

  public static async getOriginFileExts() {
    return await fetch("/api/file_exts")
      .then((res) => res.json())
      .then((data: string[]) => {
        return data.sort((a, b) =>
          a.toLowerCase() < b.toLowerCase() ? -1 : 1,
        );
      });
  }

  public static async setOriginFileExts(file_exts: string[]) {
    return await fetch(
      `/api/file_exts?file_exts=${encodeURI(JSON.stringify(file_exts))}`,
      { method: "POST" },
    );
  }

  public static async getOriginMediafiles() {
    return await fetch("/api/medialist").then((res) => res.json());
  }

  public static async refreshOriginMediafiles() {
    return await fetch("/api/medialist", { method: "POST" });
  }

  public static async getOriginLast25() {
    return await fetch("/api/last25").then((res) => res.json());
  }

  public static sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}

export { API };
