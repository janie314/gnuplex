interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

async function play() {
    return await fetch(
      `/api/play`,
      { method: "POST" },
    );
  }  
  
async function pause() {
    return await fetch(
      `/api/pause`,
      { method: "POST" },
    );
  }

async function getOriginPos() {
    return await fetch(
      `/api/pos`,
    ).then((res) => res.json()).then((res: IMPVRes) => {
      if (res.data !== undefined) {
        // @ts-ignore
        return Math.floor(res.data);
      } else {
        return 0;
      }
    });
}

async function setOriginPos(pos: number) {
    return await fetch(
      `/api/pos?pos=${pos}`,
      { method: "POST" },
    ).then((res) => res.json());
  }

export {play, pause, getOriginPos, setOriginPos};