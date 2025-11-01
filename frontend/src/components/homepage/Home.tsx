'use client'

import React, { useState } from 'react';
import { Video, Users, Brain, FileText, MessageSquare, Search, Zap, Shield, Clock, TrendingUp, Briefcase, Scale, UserCheck, ChevronRight, Moon, Sun } from 'lucide-react';

const HomePage = () => {
  const [darkMode, setDarkMode] = useState(false);

  const features = [
    {
      icon: <Brain className="w-6 h-6" />,
      title: "AI Agent Participant",
      description: "An intelligent AI joins your meetings, listens actively, and provides insights when requested."
    },
    {
      icon: <Search className="w-6 h-6" />,
      title: "Real-Time Research",
      description: "AI performs market research, competitive analysis, and information gathering during discussions."
    },
    {
      icon: <FileText className="w-6 h-6" />,
      title: "Auto Transcription & Summary",
      description: "Automatic meeting transcripts and intelligent summaries delivered to your room."
    },
    {
      icon: <MessageSquare className="w-6 h-6" />,
      title: "Contextual AI Assistant",
      description: "Chat with AI about past discussions, decisions, and key points from your room history."
    },
    {
      icon: <Users className="w-6 h-6" />,
      title: "Collaborative Rooms",
      description: "Unified workspace for video calls, messages, files, and AI insights in one place."
    },
    {
      icon: <Zap className="w-6 h-6" />,
      title: "AI-Powered Actions",
      description: "Request AI to execute tasks, gather data, and provide strategic recommendations on-demand."
    }
  ];

  const useCases = [
    {
      icon: <Briefcase className="w-8 h-8" />,
      title: "Business Development & Partnerships",
      description: "Collaborate with external partners on joint ventures, strategy, and risk assessment with AI-powered insights.",
      benefits: ["Strategic planning support", "Risk analysis", "Action item tracking"]
    },
    {
      icon: <TrendingUp className="w-8 h-8" />,
      title: "Executive Strategy Sessions",
      description: "Enhance leadership meetings with real-time market research and competitive intelligence from AI.",
      benefits: ["Market analysis", "Competitive positioning", "Data-driven decisions"]
    },
    {
      icon: <UserCheck className="w-8 h-8" />,
      title: "Sales & Account Management",
      description: "Dedicated rooms per client with conversation history, file sharing, and AI-generated meeting summaries.",
      benefits: ["Client relationship hub", "Meeting history", "Automated follow-ups"]
    },
    {
      icon: <Scale className="w-8 h-8" />,
      title: "Legal & Consulting Services",
      description: "Secure collaboration spaces for lawyers, consultants, and clients with comprehensive documentation.",
      benefits: ["Secure discussions", "Complete transcripts", "Contextual reference"]
    }
  ];

  const benefits = [
    {
      icon: <Clock className="w-6 h-6" />,
      title: "Save 5+ Hours Weekly",
      description: "Eliminate manual note-taking and post-meeting documentation"
    },
    {
      icon: <Brain className="w-6 h-6" />,
      title: "Smarter Decisions",
      description: "Access real-time insights and research during critical discussions"
    },
    {
      icon: <Shield className="w-6 h-6" />,
      title: "Never Miss Context",
      description: "AI maintains complete history and context of all room activities"
    }
  ];

  return (
    <div className={darkMode ? 'dark' : ''}>
      <div className="min-h-screen bg-gradient-to-b from-white to-gray-50 dark:from-gray-900 dark:to-gray-800 transition-colors duration-200">
        {/* Navigation */}
        {/*
        <nav className="border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              <div className="flex items-center space-x-3">
                <div className="bg-blue-600 p-2 rounded-lg">
                  <Video className="w-6 h-6 text-white" />
                </div>
                <span className="text-xl font-bold text-gray-900 dark:text-white">MeetWith.AI</span>
              </div>
              <div className="flex items-center space-x-4">
                <button
                  onClick={() => setDarkMode(!darkMode)}
                  className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-300"
                >
                  {darkMode ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
                </button>
                <button className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-colors">
                  Contact for Demo
                </button>
              </div>
            </div>
          </div>
        </nav>
        */}

        {/* Hero Section */}
        <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20 pb-16">
          <div className="text-center">
            <div className="inline-flex items-center space-x-2 bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 px-4 py-2 rounded-full mb-6">
              <Brain className="w-4 h-4" />
              <span className="text-sm font-medium">AI-Powered Video Conferencing</span>
            </div>
            <h1 className="text-5xl md:text-6xl font-bold text-gray-900 dark:text-white mb-6 leading-tight">
              Meet with Intelligence,<br />
              <span className="text-blue-600 dark:text-blue-400">Collaborate with AI</span>
            </h1>
            <p className="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-3xl mx-auto">
              The only video conferencing platform where AI participates as your intelligent team member—listening, researching, and delivering insights when you need them most.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button className="px-8 py-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg text-lg transition-colors flex items-center justify-center space-x-2">
                <span>Request a Demo</span>
                <ChevronRight className="w-5 h-5" />
              </button>
              <button className="px-8 py-4 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 text-gray-900 dark:text-white font-semibold rounded-lg text-lg border-2 border-gray-200 dark:border-gray-600 transition-colors">
                Watch Video
              </button>
            </div>
          </div>

          {/* Key Benefits */}
          <div className="grid md:grid-cols-3 gap-8 mt-20">
            {benefits.map((benefit, index) => (
              <div key={index} className="text-center">
                <div className="inline-flex items-center justify-center w-12 h-12 bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded-lg mb-4">
                  {benefit.icon}
                </div>
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">{benefit.title}</h3>
                <p className="text-gray-600 dark:text-gray-400">{benefit.description}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Features Section */}
        <section className="bg-white dark:bg-gray-900 py-20">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center mb-16">
              <h2 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
                Powerful Features for Modern Teams
              </h2>
              <p className="text-xl text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">
                Everything you need to transform your meetings from time-sinks into strategic advantage
              </p>
            </div>
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
              {features.map((feature, index) => (
                <div key={index} className="p-6 border border-gray-200 dark:border-gray-700 rounded-xl hover:shadow-lg dark:hover:shadow-blue-900/20 transition-shadow bg-white dark:bg-gray-800">
                  <div className="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded-lg flex items-center justify-center mb-4">
                    {feature.icon}
                  </div>
                  <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">{feature.title}</h3>
                  <p className="text-gray-600 dark:text-gray-400">{feature.description}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Use Cases Section */}
        <section className="py-20 bg-gradient-to-b from-gray-50 to-white dark:from-gray-800 dark:to-gray-900">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center mb-16">
              <h2 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
                Built for Enterprise Collaboration
              </h2>
              <p className="text-xl text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">
                Trusted by forward-thinking organizations across industries
              </p>
            </div>
            <div className="grid md:grid-cols-2 gap-8">
              {useCases.map((useCase, index) => (
                <div key={index} className="p-8 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 hover:shadow-xl dark:hover:shadow-blue-900/20 transition-shadow">
                  <div className="flex items-start space-x-4 mb-4">
                    <div className="w-14 h-14 bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded-lg flex items-center justify-center flex-shrink-0">
                      {useCase.icon}
                    </div>
                    <div>
                      <h3 className="text-2xl font-semibold text-gray-900 dark:text-white mb-2">{useCase.title}</h3>
                      <p className="text-gray-600 dark:text-gray-400 mb-4">{useCase.description}</p>
                    </div>
                  </div>
                  <div className="space-y-2 ml-18">
                    {useCase.benefits.map((benefit, idx) => (
                      <div key={idx} className="flex items-center space-x-2">
                        <div className="w-1.5 h-1.5 bg-blue-600 dark:bg-blue-400 rounded-full"></div>
                        <span className="text-sm text-gray-700 dark:text-gray-300">{benefit}</span>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* How It Works */}
        <section className="py-20 bg-white dark:bg-gray-900">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="text-center mb-16">
              <h2 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
                How It Works
              </h2>
            </div>
            <div className="grid md:grid-cols-3 gap-12">
              <div className="text-center">
                <div className="w-16 h-16 bg-blue-600 text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4">1</div>
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">Create a Room</h3>
                <p className="text-gray-600 dark:text-gray-400">Set up a dedicated space for your project, client, or team with one click</p>
              </div>
              <div className="text-center">
                <div className="w-16 h-16 bg-blue-600 text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4">2</div>
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">Invite Participants & AI</h3>
                <p className="text-gray-600 dark:text-gray-400">Add team members and activate your AI agent to join as an intelligent participant</p>
              </div>
              <div className="text-center">
                <div className="w-16 h-16 bg-blue-600 text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4">3</div>
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">Collaborate with Intelligence</h3>
                <p className="text-gray-600 dark:text-gray-400">Meet, discuss, and leverage AI insights—all in one unified workspace</p>
              </div>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-20 bg-gradient-to-r from-blue-600 to-blue-700 dark:from-blue-700 dark:to-blue-800">
          <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
            <h2 className="text-4xl font-bold text-white mb-6">
              Ready to Transform Your Meetings?
            </h2>
            <p className="text-xl text-blue-100 mb-8">
              Join leading enterprises using AI to make every meeting more productive and strategic
            </p>
            <button className="px-8 py-4 bg-white hover:bg-gray-100 text-blue-600 font-semibold rounded-lg text-lg transition-colors inline-flex items-center space-x-2">
              <span>Schedule Your Demo Today</span>
              <ChevronRight className="w-5 h-5" />
            </button>
          </div>
        </section>

        {/* Footer */}
        {/*
        <footer className="bg-gray-900 dark:bg-black text-gray-400 py-12">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex flex-col md:flex-row justify-between items-center">
              <div className="flex items-center space-x-3 mb-4 md:mb-0">
                <div className="bg-blue-600 p-2 rounded-lg">
                  <Video className="w-6 h-6 text-white" />
                </div>
                <span className="text-xl font-bold text-white">MeetWith.AI</span>
              </div>
              <div className="text-center md:text-right">
                <p>&copy; 2025 MeetWith.AI. All rights reserved.</p>
              </div>
            </div>
          </div>
        </footer>
        */}
      </div>
    </div>
  );
};

export default HomePage;