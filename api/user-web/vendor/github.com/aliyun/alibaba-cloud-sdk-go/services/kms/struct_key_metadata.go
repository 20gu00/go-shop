package kms

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// KeyMetadata is a nested struct in kms response
type KeyMetadata struct {
	NextRotationDate              string `json:"NextRotationDate" xml:"NextRotationDate"`
	ProtectionLevel               string `json:"ProtectionLevel" xml:"ProtectionLevel"`
	KeySpec                       string `json:"KeySpec" xml:"KeySpec"`
	DeleteDate                    string `json:"DeleteDate" xml:"DeleteDate"`
	KeyUsage                      string `json:"KeyUsage" xml:"KeyUsage"`
	AutomaticRotation             string `json:"AutomaticRotation" xml:"AutomaticRotation"`
	DKMSInstanceId                string `json:"DKMSInstanceId" xml:"DKMSInstanceId"`
	LastRotationDate              string `json:"LastRotationDate" xml:"LastRotationDate"`
	RotationInterval              string `json:"RotationInterval" xml:"RotationInterval"`
	MaterialExpireTime            string `json:"MaterialExpireTime" xml:"MaterialExpireTime"`
	PrimaryKeyVersion             string `json:"PrimaryKeyVersion" xml:"PrimaryKeyVersion"`
	DeletionProtection            string `json:"DeletionProtection" xml:"DeletionProtection"`
	Arn                           string `json:"Arn" xml:"Arn"`
	KeyState                      string `json:"KeyState" xml:"KeyState"`
	CreationDate                  string `json:"CreationDate" xml:"CreationDate"`
	Creator                       string `json:"Creator" xml:"Creator"`
	Origin                        string `json:"Origin" xml:"Origin"`
	KeyId                         string `json:"KeyId" xml:"KeyId"`
	Description                   string `json:"Description" xml:"Description"`
	DeletionProtectionDescription string `json:"DeletionProtectionDescription" xml:"DeletionProtectionDescription"`
}
