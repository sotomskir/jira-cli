#!/usr/bin/env bash
setup_git() {
  git config --global user.email "travis@travis-ci.org"
  git config --global user.name "Travis CI"
}

push_tag() {
  git tag -f nightly
  # Remove existing "origin"
  git remote rm origin
  # Add new "origin" with access token in the git URL for authentication
  git remote add origin https://sotomskir:${GH_TOKEN}@github.com/${TRAVIS_REPO_SLUG}.git > /dev/null 2>&1
  git push -f --tags
}

if [[ "${TRAVIS_PULL_REQUEST}" == "false" ]] && [[ "${TRAVIS_BRANCH}"  == "master" ]] && [[ "${TRAVIS_TAG}" == "" ]]; then
    setup_git
    push_tag
fi
