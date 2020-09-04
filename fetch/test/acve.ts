import Axios, { AxiosResponse } from "axios";

export abstract class API {
  constructor(
    protected baseUrl: string,
    protected urlParams: { key: string; preselected: string }
  ) {}

  abstract handleError(error: any);

  async LoadImageLettre(params: { lien: string }) {
    try {
      const fullUrl = this.baseUrl + "/api/imageslettre";
      const rep: AxiosResponse<any> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async PaiementHelloAsso(params: any) {
    try {
      const fullUrl = this.baseUrl + "/helloasso/nouveau_paiement";
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async LoadSondages(params: any) {
    try {
      const fullUrl =
        this.baseUrl + "/sondages/:key".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<SondageDetails[] | null> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async SearchMail(params: any) {
    try {
      const fullUrl = this.baseUrl + "/lien_espace_perso";
      const rep: AxiosResponse<SearchMailOut> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async SertMiniature(params: { "crypted-id": string; mode: string }) {
    try {
      const fullUrl = this.baseUrl + "/miniature";
      const rep: AxiosResponse<any> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async SertDocument(params: { "crypted-id": string; mode: string }) {
    try {
      const fullUrl = this.baseUrl + "/document";
      const rep: AxiosResponse<any> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async HandleUpload(params: { "crypted-id": string }, file: File) {
    try {
      const fullUrl = this.baseUrl + "/document";
      const formData = new FormData();
      formData.append("file", file, file.name);
      formData.append("crypted-id", params["crypted-id"]);
      const rep: AxiosResponse<PublicDocument> = await Axios.post(
        fullUrl,
        formData
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async HandleDelete() {
    try {
      const fullUrl = this.baseUrl + "/document";
      const rep: AxiosResponse<any> = await Axios.delete(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async FormulaireInscription() {
    try {
      const fullUrl = this.baseUrl + "/inscription";
      const rep: AxiosResponse<any> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async FormulaireInscription() {
    try {
      const fullUrl =
        this.baseUrl +
        "/inscription/:preselected".replace(
          ":preselected",
          this.urlParams.preselected
        );
      const rep: AxiosResponse<any> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async HandleLoadDataInscriptions(params: {
    preselected: string;
    preinscription: string;
  }) {
    try {
      const fullUrl = this.baseUrl + "/inscription/api";
      const rep: AxiosResponse<DataInscription> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async EnregistreInscriptionComplete(params: any) {
    try {
      const fullUrl = this.baseUrl + "/inscription/api";
      const rep: AxiosResponse<EnregistreInscriptionOut> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async RejoueInscriptionComplete(params: any) {
    try {
      const fullUrl = this.baseUrl + "/inscription/api/no-check-mail";
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async CheckMail(params: { mail: string }) {
    try {
      const fullUrl = this.baseUrl + "/inscription/api/check-mail";
      const rep: AxiosResponse<any> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async ValideMail() {
    try {
      const fullUrl = this.baseUrl + "/inscription/valide-mail";
      const rep: AxiosResponse<any> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async EnregistreInscriptionSimple(params: any) {
    try {
      const fullUrl = this.baseUrl + "/inscription/api/simple";
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async GetMetas() {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/metas".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<MetaEspacePerso> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async GetData() {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/datas".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<ContentEspacePerso> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async GetFinances() {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/finances".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<Finances> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async LoadJoomeo() {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/joomeo".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<JoomeoOutput> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async CreateDocument(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/document".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<PublicDocument> = await Axios.put(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async CreateAide(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/aide".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<Aide> = await Axios.put(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async UpdateAide(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/aide".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<Aide> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async DeleteAide() {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/aide".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<any> = await Axios.delete(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async TransfertFicheSanitaire(params: { "id-crypted": string }) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/fiche_sanitaire".replace(
          ":key",
          this.urlParams.key
        );
      const rep: AxiosResponse<string[] | null> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async UpdateFicheSanitaire(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/fiche_sanitaire".replace(
          ":key",
          this.urlParams.key
        );
      const rep: AxiosResponse<OutFicheSanitaire> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async UpdateOptionsParticipants(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/participants".replace(
          ":key",
          this.urlParams.key
        );
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async ValideTransfertFicheSanitaire(params: { target: string }) {
    try {
      const fullUrl = this.baseUrl + "/partage_fiche_sanitaire";
      const rep: AxiosResponse<any> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async MarkConnection() {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/connected".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<any> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async DownloadFacture(params: { "index-destinataire": string }) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/download/facture".replace(
          ":key",
          this.urlParams.key
        );
      const rep: AxiosResponse<any> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async DownloadAttestationPresence(params: { "index-destinataire": string }) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/download/presence".replace(
          ":key",
          this.urlParams.key
        );
      const rep: AxiosResponse<any> = await Axios.get(fullUrl, {
        params: params,
      });
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async SaveSondage(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/sondage".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<PublicSondage> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async CreeMessage(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/message".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<PseudoMessage> = await Axios.put(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async EditMessage(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/message".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<PseudoMessage> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async DeleteMessage() {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/message".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<PseudoMessage> = await Axios.delete(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async ConfirmePlaceliberee(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/espace_perso/:key/api/placeliberee".replace(
          ":key",
          this.urlParams.key
        );
      const rep: AxiosResponse<ContentEspacePerso> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async GetMembres() {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/membres".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<Membre[] | null> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async InviteOne(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/membres".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<any> = await Axios.put(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async InviteAll(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/membres".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<any> = await Axios.post(fullUrl, params);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async GetVotes() {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/votes".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<VoteAdmin[] | null> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async CreateVote(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/votes".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<VoteAdmin[] | null> = await Axios.put(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async UpdateVote(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/votes".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<VoteAdmin[] | null> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async DeleteVote() {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/votes".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<VoteCandidats[] | null> = await Axios.delete(
        fullUrl
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async ClearVote() {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/votes/clear".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<VoteAdmin[] | null> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async LockVote(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/votes/lock".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<VoteAdmin[] | null> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async ExportBilanVotes() {
    try {
      const fullUrl =
        this.baseUrl +
        "/vote_admin/:key/api/votes/export".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<any> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async PageVote() {
    try {
      const fullUrl =
        this.baseUrl + "/vote/:key".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<any> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async LoadVotes(params: any) {
    try {
      const fullUrl =
        this.baseUrl + "/vote/:key".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<VotePersonneComplet[] | null> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async EffectueVote(params: any) {
    try {
      const fullUrl =
        this.baseUrl + "/vote/:key".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<VotePersonneComplet[] | null> = await Axios.put(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async ResetVote() {
    try {
      const fullUrl =
        this.baseUrl + "/vote/:key".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<
        VotePersonneComplet[] | null
      > = await Axios.delete(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async ExportBilanPersonne() {
    try {
      const fullUrl =
        this.baseUrl + "/vote/:key/export".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<any> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async GetPasswords() {
    try {
      const fullUrl =
        this.baseUrl +
        "/passwords/:key/api".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<GetPasswordsOut> = await Axios.get(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async CreatePassword(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/passwords/:key/api".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<PublicPassword> = await Axios.put(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async UpdatePassword(params: any) {
    try {
      const fullUrl =
        this.baseUrl +
        "/passwords/:key/api".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<PublicPassword> = await Axios.post(
        fullUrl,
        params
      );
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }

  async DeletePassword() {
    try {
      const fullUrl =
        this.baseUrl +
        "/passwords/:key/api".replace(":key", this.urlParams.key);
      const rep: AxiosResponse<any> = await Axios.delete(fullUrl);
      return rep.data;
    } catch (error) {
      this.handleError(error);
    }
  }
}
