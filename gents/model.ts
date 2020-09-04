import Axios from "axios";

export class API {
  constructor(private param1: string | number) {}

  abstract handleError(error: any);

  baseUrl: string = "";

  async apiName(params: TypeIn) {
    try {
      const rep: AxiosResponse<TypeOut> = await Axios.methodName(
        this.baseUrl + url,
        { query: params }
      ); // GET, DELETE
      const rep = await Axios.methodName(this.baseUrl + url, params); // POST, PUT
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }
}
