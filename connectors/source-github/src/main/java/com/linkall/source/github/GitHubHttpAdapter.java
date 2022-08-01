package com.linkall.source.github;

import com.linkall.vance.core.Adapter2;
import io.cloudevents.CloudEvent;
import io.cloudevents.core.builder.CloudEventBuilder;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.json.JsonObject;
import org.apache.commons.lang.StringUtils;

import java.net.URI;
import java.time.*;


public class GitHubHttpAdapter implements Adapter2<HttpServerRequest, Buffer> {

    public static final CloudEventBuilder template = CloudEventBuilder.v1();

    @Override
    public CloudEvent adapt(HttpServerRequest req, Buffer buffer) {
        template.withId(req.getHeader("X-GitHub-Delivery"));
        switch (req.getHeader("X-GitHub-Event")){
            case "star": adaptStar(req, buffer);
                break;
            case "push": adaptPush(req, buffer);
                break;
            case "issues": adaptIssues(req, buffer);
                break;
            case "check_run": adaptCheckRun(req, buffer);
                break;
            case "check_suite": adaptCheckSuite(req, buffer);
                break;
            case "commit_comment": adaptCommitComment(req, buffer);
                break;
            case "content_reference": adaptContentReference(req, buffer);
                break;
            case "create": adaptCreate(req, buffer);
                break;
            case "delete": adaptDelete(req, buffer);
                break;
            case "deploy_key": adaptDeployKey(req, buffer);
                break;
            case "deployment": adaptDeployment(req, buffer);
                break;
            case "deployment_status": adaptDeploymentStatus(req, buffer);
                break;
            case "fork": adaptFork(req, buffer);
                break;
            case "github_app_authorization": adaptGitHubAppAuthorization(req, buffer);
                break;
            case "gollum": adaptGollum(req, buffer);
                break;
            case "installation": adaptInstallation(req, buffer);
                break;
            case "installation_repositories": adaptInstallationRepository(req, buffer);
                break;
            case "issue_comment": adaptIssueComment(req, buffer);
                break;
            case "label": adaptLabel(req, buffer);
                break;
            case "marketplace_purchase": adaptMarketplacePurchase(req, buffer);
                break;
            case "member": adaptMember(req, buffer);
                break;
            case "membership": adaptMemberShip(req, buffer);
                break;
            case "meta": adaptMeta(req, buffer);
                break;
            case "milestone": adaptMilestone(req, buffer);
                break;
            case "organization": adaptOrganization(req, buffer);
                break;
            case "org_block": adaptOrgBlock(req, buffer);
                break;
            case "page_build": adaptPageBuildEvent(req, buffer);
                break;
            case "project_card": adaptProjectCard(req, buffer);
                break;
            case "project_column": adaptProjectColumn(req, buffer);
                break;
            case "project": adaptProject(req, buffer);
                break;
            case "public": adaptPublic(req, buffer);
                break;
            case "pull_request": adaptPullRequest(req, buffer);
                break;
            case "pull_request_review": adaptPullRequestReview(req, buffer);
                break;
            case "pull_request_review_comment": adaptPullRequestReviewComment(req, buffer);
                break;
            case "registry_package": adaptRegistryPackage(req, buffer);
                break;
            case "release": adaptRelease(req, buffer);
                break;
            case "repository": adaptRepository(req, buffer);
                break;
            case "repository_import": adaptRepositoryImport(req, buffer);
                break;
            case "repository_vulnerability_alert": adaptRepositoryVulnerability(req, buffer);
                break;
            case "security_advisory": adaptSecurityAdvisoryEvent(req, buffer);
                break;
            case "status": adaptStatus(req, buffer);
                break;
            case "team": adaptTeam(req, buffer);
                break;
            case "team_add": adaptTeamAddEvent(req, buffer);
                break;
            case "watch": adaptWatchEvent(req, buffer);
                break;
            default: adaptOtherEvent(req, buffer);
        }
        JsonObject data = new JsonObject();
        String contentType = req.getHeader("content-type");
        if(null != contentType && contentType.equals("application/json")){
            JsonObject body = buffer.toJsonObject();
            data.put("body", body);
        }else{
            String myData = new String(buffer.getBytes());
            JsonObject body = new JsonObject();
            data.put("data", myData);
            data.put("body", body);
        }
        template.withData(data.toBuffer().getBytes());
        return template.build();
    }

    private void adaptOtherEvent(HttpServerRequest req, Buffer buffer){
        //解析时间
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withDataContentType("application/json")
                .withTime(eventTime);
    }

    private void adaptStar(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String html_url = repository.getString("url");
        URI uri = URI.create(html_url);
        String action = payLoad.getString("action");
        String time = repository.getString("starred_at");
        //解析时间
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.star."+action)
                .withDataContentType("application/json")
                .withTime(eventTime);
    }

    private void adaptPush(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String html_url = repository.getString("url");
        URI uri = URI.create(html_url);
        String subject = payLoad.getString("ref");
        String time = repository.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.push")
                .withDataContentType("application/json")
                .withTime(eventTime)
                .withSubject(subject);
    }

    private void adaptIssues(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String html_url = repository.getString("url");
        URI uri = URI.create(html_url);
        String action = payLoad.getString("action");
        JsonObject issue = payLoad.getJsonObject("issue");
        String number = issue.getString("number");
        String time = issue.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.issue."+action)
                .withSubject(number)
                .withDataContentType("application/json")
                .withTime(eventTime);
    }

    private void adaptCheckRun(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String html_url = repository.getString("url");
        URI uri = URI.create(html_url);
        String action = payLoad.getString("action");
        JsonObject check_run = payLoad.getJsonObject("check_run");
        String subject = check_run.getString("id");
        String time = check_run.getString("completed_at");
        if(StringUtils.isBlank(time)){
            time = check_run.getString("started_at");
        }

        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.check_run."+action)
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }
    private void adaptCheckSuite(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String html_url = repository.getString("url");
        URI uri = URI.create(html_url);
        String action = payLoad.getString("action");
        JsonObject check_run = payLoad.getJsonObject("check_suite");
        String subject = check_run.getString("id");
        String time = check_run.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.check_suite."+action)
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }
    private void adaptCommitComment(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject comment = payLoad.getJsonObject("comment");
        String source = comment.getString("url")+"/"+comment.getString("comment_id");
        URI uri = URI.create(source);
        String action = payLoad.getString("action");
        String subject = comment.getString("id");
        String time = comment.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.commit_comment."+action)
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }
    private void adaptContentReference(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        String action = payLoad.getString("action");
        JsonObject content_reference = payLoad.getJsonObject("content_reference");
        String subject = content_reference.getString("id");
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.content_reference."+action)
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }
    private void adaptCreate(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        String ref_type = payLoad.getString("ref_type");
        String subject = payLoad.getString("ref");
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.create."+ref_type)
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }
    private void adaptDelete(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        String ref_type = payLoad.getString("ref_type");
        String subject = payLoad.getString("ref");
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.delete."+ref_type)
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }
    private void adaptDeployKey(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        String action = payLoad.getString("action");
        JsonObject key = payLoad.getJsonObject("key");
        String subject = key.getString("id");
        String time = key.getString("deleted_at");
        if(StringUtils.isBlank(time)){
            time = key.getString("created_at");
        }
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.deploy_key."+action)
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }
    private void adaptDeployment(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        JsonObject deployment = payLoad.getJsonObject("deployment");
        String subject = deployment.getString("id");
        String time = deployment.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.deployment")
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }

    private void adaptDeploymentStatus(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject deployment = payLoad.getJsonObject("deployment");
        String source = deployment.getString("url");
        URI uri = URI.create(source);
        JsonObject deployment_status = payLoad.getJsonObject("deployment_status");
        String subject = deployment_status.getString("url");
        String time = deployment_status.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.deployment_status."+deployment_status.getString("state"))
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }

    private void adaptFork(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        JsonObject forkee = payLoad.getJsonObject("forkee");
        String subject = forkee.getString("url");
        String time = forkee.getString("created_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.fork")
                .withSubject(subject)
                .withTime(eventTime)
                .withDataContentType("application/json");
    }

    private void adaptGitHubAppAuthorization(HttpServerRequest req, Buffer buffer){
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject sender = payLoad.getJsonObject("sender");
        String source = sender.getString("url");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.github_app_authorization")
                .withTime(eventTime)
                .withDataContentType("application/json");
    }

    private void adaptGollum(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject pages = payLoad.getJsonObject("pages");
        String action = pages.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.gollum." + action)
                .withTime(eventTime)
                .withSubject(pages.getString("page_name"))
                .withDataContentType("application/json");
    }

    private void adaptInstallation(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject installation = payLoad.getJsonObject("installation");
        JsonObject account = installation.getJsonObject("account");
        String source = account.getString("url");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        Long time = Long.parseLong(installation.getString("updated_at"));
        OffsetDateTime eventTime = getTimeFromSecondTimestamp(time);
        template.withSource(uri)
                .withType("com.github.installation." + action)
                .withTime(eventTime)
                .withSubject(installation.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptInstallationRepository(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject installation = payLoad.getJsonObject("installation");
        JsonObject account = installation.getJsonObject("account");
        String source = account.getString("url");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        Long time = Long.parseLong(installation.getString("updated_at"));
        OffsetDateTime eventTime = getTimeFromSecondTimestamp(time);
        template.withSource(uri)
                .withType("com.github.installation_repository." + action)
                .withTime(eventTime)
                .withSubject(installation.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptIssueComment(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject issue = payLoad.getJsonObject("issue");
        String source = issue.getString("url");
        JsonObject comment = payLoad.getJsonObject("comment");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(comment.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.issue_comment." + action)
                .withTime(eventTime)
                .withSubject(comment.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptLabel(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject label = payLoad.getJsonObject("label");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.label." + action)
                .withTime(eventTime)
                .withSubject(label.getString("name"))
                .withDataContentType("application/json");
    }

    private void adaptMarketplacePurchase(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject sender = payLoad.getJsonObject("sender");
        String source = sender.getString("url").replace("/username", "");
        JsonObject marketplace_purchase = payLoad.getJsonObject("label");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(payLoad.getString("effective_date"));
        template.withSource(uri)
                .withType("com.github.marketplace_purchase." + action)
                .withTime(eventTime)
                .withSubject(marketplace_purchase.getString("name"))
                .withDataContentType("application/json");
    }

    private void adaptMember(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject member = payLoad.getJsonObject("member");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.member." + action)
                .withTime(eventTime)
                .withSubject(member.getString("login"))
                .withDataContentType("application/json");
    }

    private void adaptMemberShip(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject team = payLoad.getJsonObject("team");
        String source = team.getString("url");
        JsonObject member = payLoad.getJsonObject("member");
        String scope = payLoad.getString("scope");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.membership." + scope + action)
                .withTime(eventTime)
                .withSubject(member.getString("login"))
                .withDataContentType("application/json");
    }

    private void adaptMeta(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject hook = payLoad.getJsonObject("hook");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(hook.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.meta." + action)
                .withTime(eventTime)
                .withSubject(payLoad.getString("hook_id"))
                .withDataContentType("application/json");
    }

    private void adaptMilestone(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject milestone = payLoad.getJsonObject("milestone");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(milestone.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.milestone." + action)
                .withTime(eventTime)
                .withSubject(milestone.getString("number"))
                .withDataContentType("application/json");
    }

    private void adaptOrganization(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject organization = payLoad.getJsonObject("organization");
        String source = organization.getString("url");
        JsonObject membership = payLoad.getJsonObject("membership");
        JsonObject user = membership.getJsonObject("user");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.organization." + action)
                .withTime(eventTime)
                .withSubject(user.getString("login"))
                .withDataContentType("application/json");
    }

    private void adaptOrgBlock(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject organization = payLoad.getJsonObject("organization");
        String source = organization.getString("url");
        JsonObject blocked_user = payLoad.getJsonObject("blocked_user");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.org_block." + action)
                .withTime(eventTime)
                .withSubject(blocked_user.getString("login"))
                .withDataContentType("application/json");
    }

    private void adaptPageBuildEvent(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject build = payLoad.getJsonObject("build");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(payLoad.getJsonObject("pusher").getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.page_build")
                .withTime(eventTime)
                .withSubject(build.getString("url"))
                .withDataContentType("application/json");
    }

    private void adaptProjectCard(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject project_card = payLoad.getJsonObject("project_card");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(project_card.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.project_card."+action)
                .withTime(eventTime)
                .withSubject(project_card.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptProjectColumn(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject project_column = payLoad.getJsonObject("project_column");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(project_column.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.project_column."+action)
                .withTime(eventTime)
                .withSubject(project_column.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptProject(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject project = payLoad.getJsonObject("project");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(project.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.project."+action)
                .withTime(eventTime)
                .withSubject(project.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptPublic(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        JsonObject owner = repository.getJsonObject("owner");
        String source = owner.getString("url");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(repository.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.public")
                .withTime(eventTime)
                .withSubject(repository.getString("name"))
                .withDataContentType("application/json");
    }

    private void adaptPullRequest(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(repository.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.pull_request." + action)
                .withTime(eventTime)
                .withSubject(payLoad.getString("number"))
                .withDataContentType("application/json");
    }

    private void adaptPullRequestReview(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject pull_request = payLoad.getJsonObject("pull_request");
        String source = pull_request.getString("url");
        JsonObject review = payLoad.getJsonObject("review");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(review.getString("submitted_at"));
        template.withSource(uri)
                .withType("com.github.pull_request_review." + action)
                .withTime(eventTime)
                .withSubject(review.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptPullRequestReviewComment(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject pull_request = payLoad.getJsonObject("pull_request");
        String source = pull_request.getString("url");
        JsonObject comment = payLoad.getJsonObject("comment");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(pull_request.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.pull_request_review_comment." + action)
                .withTime(eventTime)
                .withSubject(comment.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptRegistryPackage(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject registry_package = payLoad.getJsonObject("registry_package");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(registry_package.getString("updated_at"));
        template.withSource(uri)
                .withType("com.github.registry_package." + action)
                .withTime(eventTime)
                .withSubject(registry_package.getString("html_url"))
                .withDataContentType("application/json");
    }

    private void adaptRelease(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        JsonObject release = payLoad.getJsonObject("release");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(release.getString(action.equals("created_at") ? "created_at" : "published_at"));
        template.withSource(uri)
                .withType("com.github.release." + action)
                .withTime(eventTime)
                .withSubject(release.getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptRepository(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        JsonObject owner = repository.getJsonObject("owner");
        String source = owner.getString("url");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(repository.getString("updated_at" ));
        template.withSource(uri)
                .withType("com.github.repository." + action)
                .withTime(eventTime)
                .withSubject(repository.getString("name"))
                .withDataContentType("application/json");
    }

    private void adaptRepositoryImport(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        JsonObject owner = repository.getJsonObject("owner");
        String source = owner.getString("url");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getTime(repository.getString("updated_at" ));
        template.withSource(uri)
                .withType("com.github.repository_import")
                .withTime(eventTime)
                .withSubject(repository.getString("name"))
                .withDataContentType("application/json");
    }

    private void adaptRepositoryVulnerability(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        String action = payLoad.getString("action");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.repository_vulnerability_alert."+action)
                .withTime(eventTime)
                .withSubject(payLoad.getJsonObject("alert").getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptSecurityAdvisoryEvent(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject security_advisory = payLoad.getJsonObject("security_advisory");
        String action = payLoad.getString("action");
        URI uri = URI.create("github.com");
        String time = security_advisory.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.security_advisory."+action)
                .withTime(eventTime)
                .withSubject(security_advisory.getString("ghsa_id"))
                .withDataContentType("application/json");
    }

    private void adaptStatus(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        String time = payLoad.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.status."+payLoad.getString("status"))
                .withTime(eventTime)
                .withSubject(payLoad.getString("sha"))
                .withDataContentType("application/json");
    }

    private void adaptTeam(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        String time = payLoad.getString("updated_at");
        OffsetDateTime eventTime = getTime(time);
        template.withSource(uri)
                .withType("com.github.team."+payLoad.getString("action"))
                .withTime(eventTime)
                .withSubject(payLoad.getJsonObject("team").getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptTeamAddEvent(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.team_add."+payLoad.getString("action"))
                .withTime(eventTime)
                .withSubject(payLoad.getJsonObject("team").getString("id"))
                .withDataContentType("application/json");
    }

    private void adaptWatchEvent(HttpServerRequest req, Buffer buffer) {
        JsonObject payLoad = buffer.toJsonObject();
        JsonObject repository = payLoad.getJsonObject("repository");
        String source = repository.getString("url");
        URI uri = URI.create(source);
        OffsetDateTime eventTime = getZeroTime(LocalDateTime.now());
        template.withSource(uri)
                .withType("com.github.watch."+payLoad.getString("action"))
                .withTime(eventTime)
                .withDataContentType("application/json");
    }

    private  OffsetDateTime getTime(String time){
        if(StringUtils.isBlank(time)){
            return getZeroTime(LocalDateTime.now());
        }else{
            time = time.substring(0, time.length() - 1);
            LocalDateTime localDateTime = LocalDateTime.parse(time);
            return OffsetDateTime.of(localDateTime, ZoneOffset.UTC);
        }
    }

    private OffsetDateTime getZeroTime(LocalDateTime time){
        LocalDateTime dt = LocalDateTime.now(ZoneId.of("Z"));
        Duration duration = Duration.between(time, dt);
        OffsetDateTime time2 = OffsetDateTime.of(time, ZoneOffset.UTC).plus(duration);
        return time2;
    }

    private OffsetDateTime getTimeFromSecondTimestamp(long timestamp){
        return Instant.ofEpochSecond(timestamp).atOffset(ZoneOffset.UTC);
    }
}
