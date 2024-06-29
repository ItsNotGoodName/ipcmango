import { A } from "@solidjs/router";
import {
  BreadcrumbsRoot,
  BreadcrumbsItem,
  BreadcrumbsLink,
  BreadcrumbsSeparator,
} from "~/ui/Breadcrumbs";
import { LayoutNormal } from "~/ui/Layout";
import { PageTitle } from "~/ui/Page";

export default function () {
  return (
    <LayoutNormal class="max-w-4xl">
      <PageTitle>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            <BreadcrumbsLink as={A} href="/events">
              Events
            </BreadcrumbsLink>
            <BreadcrumbsSeparator />
          </BreadcrumbsItem>
          <BreadcrumbsItem>Rules</BreadcrumbsItem>
        </BreadcrumbsRoot>
      </PageTitle>
    </LayoutNormal>
  );
}
