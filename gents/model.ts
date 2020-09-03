import Axios from "axios"

export class API {

    abstract handleError(error: any) {}

    async apiName(params: TypeIn) {
        try {
            const rep:<AxiosResponse<TypeOut>> = await Axios.methodName(url, {query: params}) // GET, DELETE
            const rep = await Axios.methodName(url, params) // POST, PUT
            return rep.data
        } catch (error) {
            this.handleError(error)
        }
    }
}