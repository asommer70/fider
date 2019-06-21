import "./PrivacySettings.page.scss";

import React from "react";
import { Toggle, Select, SelectOption, Form } from "@fider/components/common";
import { actions, notify, Fider } from "@fider/services";
import { AdminBasePage } from "@fider/pages/Administration/components/AdminBasePage";
import { FaKey } from "react-icons/fa";

interface PrivacySettingsPageState {
  isPrivate: boolean;
  createPosts: string;
}

export default class PrivacySettingsPage extends AdminBasePage<{}, PrivacySettingsPageState> {
  public id = "p-admin-privacy";
  public name = "privacy";
  public icon = FaKey;
  public title = "Privacy";
  public subtitle = "Manage your site privacy";

  constructor(props: {}) {
    super(props);

    this.state = {
      isPrivate: Fider.session.tenant.isPrivate,
      createPosts: Fider.session.tenant.createPosts.toString()
    };
  }

  private toggle = async (active: boolean) => {
    this.setState(
      state => ({
        isPrivate: active
      }),
      async () => {
        const response = await actions.updateTenantPrivacy(this.state.isPrivate);
        if (response.ok) {
          notify.success("Your privacy settings have been saved.");
        }
      }
    );
  };

  private setCreatePosts = (opt?: SelectOption) => {
    if (opt) {
      this.setState(
        state => ({
          createPosts: opt.value
        }),
        async () => {
          const response = await actions.updateTenantCreatePosts(parseInt(this.state.createPosts));
          if (response.ok) {
            notify.success("Your create post settings have been saved.");
          }
        }
      );
    } 
  };

  public content() {
    const options = [
      {'value': 3,'title': 'Administrator'},
      {'value': 2,'title': 'Collaborator'},
      {'value': 1,'title': 'Visitor'},
    ].map(s => ({
      value: s.value.toString(),
      label: s.title
    }));

    return (
      <div>
        <Form>
          <div className="c-form-field">
            <label htmlFor="private">Private site</label>
            <Toggle disabled={!Fider.session.user.isAdministrator} active={this.state.isPrivate} onToggle={this.toggle} />
            <p className="info">
              A private site prevents unauthenticated users from viewing or interacting with its content. <br /> If
              enabled, only already registered and invited users will be able to sign in to this site.
            </p>
          </div>
        </Form>
        <br /><br />
        <Form>
          <div className="c-form-field">
            <Select
              field="createPosts"
              label="Create Posts"
              defaultValue={'1'}
              options={options}
              onChange={this.setCreatePosts}
            />
            <p className="info">
              All roles selected and above will have access to create posts.<br />
              For example, choosing "Collaborator" will allow both Collaborators and Administrators the ability to create posts.
            </p>
          </div>
        </Form>
      </div>
    );
  }
}
