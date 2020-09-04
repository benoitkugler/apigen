import Axios, { AxiosResponse } from "./axios";

export abstract class API {
  constructor(
    protected baseUrl: string,
    protected urlParams: { class_: string; default_: string; param: string }
  ) {}

  abstract handleError(error: any);

  async handle1() {
    try {
      const fullUrl = this.baseUrl + "const_local_url";
      const rep: AxiosResponse<string> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async handleExt(params: any) {
    try {
      const fullUrl = this.baseUrl + "/const_url_from_inner_package/";
      const rep: AxiosResponse<{
        [key: string]: number[] | null;
      } | null> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async handler2(params: any) {
    try {
      const fullUrl = this.baseUrl + "/const_url_from_inner_package/endpoint";
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async handler3(params: any) {
    try {
      const fullUrl = this.baseUrl + "host/const_url_from_inner_package/";
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async handler4(params: any) {
    try {
      const fullUrl = this.baseUrl + "hostendpoint";
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async handler5(params: any) {
    try {
      const fullUrl = this.baseUrl + "/string_litteral";
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async handler6(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/with_param/:param".replace(":param", this.urlParams.param);
      const rep: AxiosResponse<any> = await Axios.put(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async handler7() {
    try {
      const fullUrl =
        this.baseUrl +
        "/special_param_value/:class/route".replace(
          ":class",
          this.urlParams.class_
        );
      const rep: AxiosResponse<any> = await Axios.delete(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async handler8(params: {
    query_param1: string | number;
    query_param2: string | number;
  }) {
    try {
      const fullUrl =
        this.baseUrl +
        "/special_param_value/:default/route".replace(
          ":default",
          this.urlParams.default_
        );
      const rep: AxiosResponse<number> = await Axios.delete(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }
}
