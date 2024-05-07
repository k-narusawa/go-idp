type PasskeyResponse = {
  keys: PasskeyResponseItem[];
};

type PasskeyResponseItem = {
  id: string;
  aaguid: string;
  key_name: string;
};
