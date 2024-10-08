package nullplatform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const APPROVAL_ACTION_PATH = "/approval/action"

type ApprovalAction struct {
	Id              int               `json:"id,omitempty"`
	Nrn             string            `json:"nrn,omitempty"`
	Entity          string            `json:"entity,omitempty"`
	Action          string            `json:"action,omitempty"`
	Dimensions      map[string]string `json:"dimensions,omitempty"`
	OnPolicySuccess string            `json:"on_policy_success,omitempty"`
	OnPolicyFail    string            `json:"on_policy_fail,omitempty"`
	Status          string            `json:"status,omitempty"`
	Policies        []*ApprovalPolicy `json:"policies,omitempty"`
}

func (c *NullClient) CreateApprovalAction(action *ApprovalAction) (*ApprovalAction, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(*action)

	if err != nil {
		return nil, err
	}

	res, err := c.MakeRequest("POST", APPROVAL_ACTION_PATH, &buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var nErr NullErrors
		if err := json.NewDecoder(res.Body).Decode(&nErr); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		defer res.Body.Close()
		return nil, fmt.Errorf("error creating approval action resource, got status code: %d, %s", res.StatusCode, nErr.Message)
	}

	actionRes := &ApprovalAction{}
	derr := json.NewDecoder(res.Body).Decode(actionRes)

	if derr != nil {
		return nil, derr
	}

	return actionRes, nil
}

func (c *NullClient) PatchApprovalAction(approvalActionId string, action *ApprovalAction) error {
	path := fmt.Sprintf("%s/%s", APPROVAL_ACTION_PATH, approvalActionId)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(*action)

	if err != nil {
		return err
	}

	res, err := c.MakeRequest("PATCH", path, &buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusNoContent) {
		return fmt.Errorf("error patching approval action resource, got %d", res.StatusCode)
	}

	return nil
}

func (c *NullClient) GetApprovalAction(approvalActionId string) (*ApprovalAction, error) {
	path := fmt.Sprintf("%s/%s", APPROVAL_ACTION_PATH, approvalActionId)

	res, err := c.MakeRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	action := &ApprovalAction{}
	derr := json.NewDecoder(res.Body).Decode(action)

	if derr != nil {
		return nil, derr
	}

	if action.Status == "deleted" {
		return action, fmt.Errorf("error getting approval action resource, the status is %s", action.Status)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting approval action resource, got %d for %s", res.StatusCode, approvalActionId)
	}

	return action, nil
}

func (c *NullClient) DeleteApprovalAction(approvalActionId string) error {
	path := fmt.Sprintf("%s/%s", APPROVAL_ACTION_PATH, approvalActionId)

	res, err := c.MakeRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusNoContent) {
		return fmt.Errorf("error deleting approval action resource, got %d", res.StatusCode)
	}

	return nil
}

func (c *NullClient) AssociatePolicyWithAction(approvalActionId, approvalPolicyID string) error {
	var buf bytes.Buffer
	path := fmt.Sprintf("%s/%s/policy", APPROVAL_ACTION_PATH, approvalActionId)

	policy := map[string]string{
		"policy_id": approvalPolicyID,
	}
	err := json.NewEncoder(&buf).Encode(policy)

	if err != nil {
		return err
	}

	res, err := c.MakeRequest("POST", path, &buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusCreated) {
		var nErr NullErrors
		if err := json.NewDecoder(res.Body).Decode(&nErr); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		defer res.Body.Close()
		return fmt.Errorf("error associating approval policy with action, got status code: %d, %s", res.StatusCode, nErr.Message)
	}

	return nil
}

func (c *NullClient) DisassociatePolicyFromAction(approvalActionId, approvalPolicyID string) error {
	path := fmt.Sprintf("%s/%s/policy/%s", APPROVAL_ACTION_PATH, approvalActionId, approvalPolicyID)

	res, err := c.MakeRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if (res.StatusCode != http.StatusOK) && (res.StatusCode != http.StatusNoContent) {
		return fmt.Errorf("error deleting approval action and policy association, got %d", res.StatusCode)
	}

	return nil
}
