// Copyright 2018 The Chubao Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package master

import (
	"encoding/json"
	"fmt"

	"github.com/chubaofs/chubaofs/proto"
	"github.com/chubaofs/chubaofs/util/errors"
	"github.com/chubaofs/chubaofs/util/log"
)

func (u *User) submit(metadata *RaftCmd) (err error) {
	cmd, err := metadata.Marshal()
	if err != nil {
		return errors.New(err.Error())
	}
	if _, err = u.partition.Submit(cmd); err != nil {
		msg := fmt.Sprintf("action[user_submit] err:%v", err.Error())
		return errors.New(msg)
	}
	return
}

// key=#ak#accesskey,value = akPolicy
func (u *User) syncAddAKPolicy(akPolicy *proto.AKPolicy) (err error) {
	return u.syncPutAKPolicy(opSyncAddAKPolicy, akPolicy)
}

func (u *User) syncDeleteAKPolicy(akPolicy *proto.AKPolicy) (err error) {
	return u.syncPutAKPolicy(opSyncDeleteAKPolicy, akPolicy)
}

func (u *User) syncUpdateAKPolicy(akPolicy *proto.AKPolicy) (err error) {
	return u.syncPutAKPolicy(opSyncUpdateAKPolicy, akPolicy)
}

func (u *User) syncPutAKPolicy(opType uint32, akPolicy *proto.AKPolicy) (err error) {
	userInfo := new(RaftCmd)
	userInfo.Op = opType
	userInfo.K = akPrefix + akPolicy.AccessKey
	userInfo.V, err = json.Marshal(akPolicy)
	if err != nil {
		return errors.New(err.Error())
	}
	return u.submit(userInfo)
}

// key=#user#userid,value = akPolicy
func (u *User) syncAddUserAK(userAK *proto.UserAK) (err error) {
	return u.syncPutUserAK(opSyncAddUserAK, userAK)
}

func (u *User) syncDeleteUserAK(userAK *proto.UserAK) (err error) {
	return u.syncPutUserAK(opSyncDeleteUserAK, userAK)
}

func (u *User) syncPutUserAK(opType uint32, userAK *proto.UserAK) (err error) {
	userInfo := new(RaftCmd)
	userInfo.Op = opType
	userInfo.K = userPrefix + userAK.UserID
	userInfo.V, err = json.Marshal(userAK)
	if err != nil {
		return errors.New(err.Error())
	}
	return u.submit(userInfo)
}

func (u *User) syncAddVolAK(volAK *proto.VolAK) (err error) {
	return u.syncPutVolAK(opSyncAddVolAK, volAK)
}

func (u *User) syncDeleteVolAK(volAK *proto.VolAK) (err error) {
	return u.syncPutVolAK(opSyncDeleteVolAK, volAK)
}

func (u *User) syncUpdateVolAK(volAK *proto.VolAK) (err error) {
	return u.syncPutVolAK(opSyncUpdateVolAK, volAK)
}

func (u *User) syncPutVolAK(opType uint32, volAK *proto.VolAK) (err error) {
	userInfo := new(RaftCmd)
	userInfo.Op = opType
	userInfo.K = volAKPrefix + volAK.Vol
	userInfo.V, err = json.Marshal(volAK)
	if err != nil {
		return errors.New(err.Error())
	}
	return u.submit(userInfo)
}

func (u *User) loadAKStore() (err error) {
	result, err := u.fsm.store.SeekForPrefix([]byte(akPrefix))
	if err != nil {
		err = fmt.Errorf("action[loadAccessKeyInfo], err: %v", err.Error())
		return err
	}
	for _, value := range result {
		aks := &proto.AKPolicy{}
		if err = json.Unmarshal(value, aks); err != nil {
			err = fmt.Errorf("action[loadAccessKeyInfo], unmarshal err: %v", err.Error())
			return err
		}
		u.akStore.Store(aks.AccessKey, aks)
		log.LogInfof("action[loadAccessKeyInfo], ak[%v]", aks.AccessKey)
	}
	return
}

func (u *User) loadUserAK() (err error) {
	result, err := u.fsm.store.SeekForPrefix([]byte(userPrefix))
	if err != nil {
		err = fmt.Errorf("action[loadUserAK], err: %v", err.Error())
		return err
	}
	for _, value := range result {
		user := &proto.UserAK{}
		if err = json.Unmarshal(value, user); err != nil {
			err = fmt.Errorf("action[loadUserAK], unmarshal err: %v", err.Error())
			return err
		}
		u.userAk.Store(user.UserID, user)
		log.LogInfof("action[loadUserAK], userID[%v]", user.UserID)
	}
	return
}

func (u *User) loadVolAKs() (err error) {
	result, err := u.fsm.store.SeekForPrefix([]byte(volAKPrefix))
	if err != nil {
		err = fmt.Errorf("action[loadVolAKs], err: %v", err.Error())
		return err
	}
	for _, value := range result {
		volAK := &proto.VolAK{}
		if err = json.Unmarshal(value, volAK); err != nil {
			err = fmt.Errorf("action[loadVolAKs], unmarshal err: %v", err.Error())
			return err
		}
		u.volAKs.Store(volAK.Vol, volAK)
		log.LogInfof("action[loadVolAKs], vol[%v]", volAK.Vol)
	}
	return
}
