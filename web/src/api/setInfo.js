import service from "@/utils/request"
const BASE_URL = `/user/api/v1`

export const setPlatformInfo = (type, data, config) => {
    return service({
        url: `${BASE_URL}/custom/${type}`,
        method: "post",
        data,
        config
    });
}