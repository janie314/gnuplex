interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

class APICall {
  /**
   * @returns GNUPlex version string.
   */
  public static async version(): Promise<string> {
    return await fetch(
      `/api/version`,
    ).then((res) => res.json());
  }

  /**
   * @returns Whether the video player is paused.
   */
  public static async paused(): Promise<boolean | null> {
    return await fetch(
      `/api/paused`,
    ).then((res) => res.json());
  }

  /**
   * Toggles the video's play/pause status.
   *
   * @returns Whether the video player is paused following the toggle operation.
   */
  public static async toggle(): Promise<boolean> {
    return await fetch(
      `/api/toggle`,
      { method: "POST" },
    ).then((res) => res.json());
  }

  /**
   * Increments the video's position.
   *
   * @param inc How much to increment the position (seconds, can be positive or negative)
   * @returns The video player's position (seconds) following the toggle inc operation.
   */
  public static async incPos(inc: number): Promise<number | null> {
    return await fetch(
      `/api/incpos?inc=${inc}`,
      { method: "POST" },
    ).then((res) => res.json());
  }

  /**
   * @returns The video player's position (seconds).
   */
  public static async pos(): Promise<number | null> {
    return await fetch(
      `/api/pos`,
    ).then((res) => res.json());
  }

  /**
   * Sets the video's position.
   *
   * @param The position the video should seek to (seconds).
   */
  public static async setPos(pos: number) {
    return await fetch(
      `/api/pos?pos=${pos}`,
      { method: "POST" },
    ).then((res) => res.json());
  }

  /**
   * @returns The video's volume (percentage, 0 - 100+).
   */
  public static async vol(): Promise<number | null> {
    return await fetch(
      `/api/vol`,
    ).then((res) => res.json());
  }

  /**
   * @param vol The video's volume (percentage, 0 - 100+).
   */
  public static async setVol(vol: number) {
    return await fetch(
      `/api/vol?vol=${vol}`,
      { method: "POST" },
    ).then((res) => res.json());
  }

  public static async getOriginMedia() {
    return await fetch(
      `/api/media`,
    ).then((res) => res.json()).then((res: IMPVRes) => {
      if (res.data !== undefined) {
        return res.data as string;
      } else {
        return "";
      }
    });
  }

  public static async setOriginMedia(mediafile: string) {
    return await fetch(
      `/api/media?mediafile=${encodeURI(mediafile)}`,
      { method: "POST" },
    );
  }

  public static async getOriginMediadirs() {
    return await fetch(
      `/api/mediadirs`,
    ).then((res) => res.json()).then((data: string[]) => {
      return data.sort((a, b) => a.toLowerCase() < b.toLowerCase() ? -1 : 1);
    });
  }

  public static async setOriginMediadirs(mediadirs: string[]) {
    return await fetch(
      `/api/mediadirs?mediadirs=${encodeURI(JSON.stringify(mediadirs))}`,
      { method: "POST" },
    );
  }

  public static async getOriginFileExts() {
    return await fetch(
      `/api/file_exts`,
    ).then((res) => res.json()).then((data: string[]) => {
      return data.sort((a, b) => a.toLowerCase() < b.toLowerCase() ? -1 : 1);
    });
  }

  public static async setOriginFileExts(file_exts: string[]) {
    return await fetch(
      `/api/file_exts?file_exts=${encodeURI(JSON.stringify(file_exts))}`,
      { method: "POST" },
    );
  }

  public static async getOriginMediafiles() {
    return await fetch(`/api/medialist`).then((res) => res.json());
  }

  public static async refreshOriginMediafiles() {
    return await fetch(`/api/medialist`, { method: "POST" });
  }

  public static async getOriginLast25() {
    return await fetch(`/api/last25`).then((res) => res.json());
  }

  public static sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}

export { APICall };
